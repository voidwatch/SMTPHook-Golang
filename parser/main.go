package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type ParsedEmail struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Date    string `json:"date"`
	Body    string `json:"body"`
}

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Service   string      `json:"service"`
	Event     string      `json:"event"`
	Data      ParsedEmail `json:"data"`
}

func parseEmail(input string) ParsedEmail {
	msg, err := mail.ReadMessage(strings.NewReader(input))
	if err != nil {
		log.Printf("⚠️ Failed to parse email with net/mail, falling back to manual parse: %v", err)
		return parseEmailFallback(input)
	}

	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("⚠️ Failed to read email body: %v", err)
	}

	return ParsedEmail{
		Subject: msg.Header.Get("Subject"),
		From:    msg.Header.Get("From"),
		To:      msg.Header.Get("To"),
		Date:    msg.Header.Get("Date"),
		Body:    strings.TrimSpace(string(body)),
	}
}

// parseEmailFallback is a simple line-based parser used when net/mail fails.
func parseEmailFallback(input string) ParsedEmail {
	var subject, from, to, date string
	var bodyBuilder strings.Builder

	lines := strings.SplitAfter(input, "\n")
	readBody := false
	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		if !readBody && line == "" {
			readBody = true
			continue
		}
		if !readBody {
			if strings.HasPrefix(line, "Subject:") {
				subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
			} else if strings.HasPrefix(line, "From:") {
				from = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
			} else if strings.HasPrefix(line, "To:") {
				to = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
			} else if strings.HasPrefix(line, "Date:") {
				date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
			}
		} else {
			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
		}
	}

	return ParsedEmail{
		Subject: subject,
		From:    from,
		To:      to,
		Date:    date,
		Body:    strings.TrimSpace(bodyBuilder.String()),
	}
}

func postWithRetry(url string, data []byte, maxRetries int, initialDelay time.Duration) error {
	delay := initialDelay
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		log.Printf("Attempt %d/%d failed. Retrying in %s...", attempt, maxRetries, delay)
		time.Sleep(delay)
		delay *= 2 // exponential backoff
	}
	return fmt.Errorf("all %d attempts to POST failed", maxRetries)
}

func startHealthServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	go func() {
		log.Printf("Health check endpoint running at :%s/health", port)
		if err := http.ListenAndServe(net.JoinHostPort("", port), mux); err != nil {
			log.Fatalf("Health check server error: %v", err)
		}
	}()
}

func main() {
	_ = os.MkdirAll("logs", 0755)

	_ = godotenv.Load()

	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		logFilePath = "logs/parser.log"
	}
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Could not open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	inputFile := os.Getenv("EMAIL_INPUT_FILE")
	if inputFile == "" {
		log.Fatal("EMAIL_INPUT_FILE environment variable is required")
	}
	webhookURL := os.Getenv("WEBHOOK_URL")
	pollInterval := os.Getenv("POLL_INTERVAL")
	if pollInterval == "" {
		pollInterval = "10"
	}
	intervalSeconds, err := strconv.Atoi(pollInterval)
	if err != nil || intervalSeconds <= 0 {
		log.Printf("Invalid POLL_INTERVAL %q, defaulting to 10s", pollInterval)
		intervalSeconds = 10
	}
	intervalDur := time.Duration(intervalSeconds) * time.Second

	healthPort := os.Getenv("HEALTH_PORT")
	if healthPort == "" {
		healthPort = "4010"
	}
	startHealthServer(healthPort)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Parser service started. Watching for email input every", intervalDur)

loop:
	for {
		select {
		case <-stop:
			log.Println("Shutting down parser...")
			break loop
		default:
			data, err := os.ReadFile(inputFile)
			if err != nil {
				if !os.IsNotExist(err) {
					log.Printf("Failed to read input file: %v", err)
				}
				time.Sleep(intervalDur)
				continue
			}
			if len(data) == 0 {
				time.Sleep(intervalDur)
				continue
			}

			// Atomically replace file content to avoid race condition
			if err := os.WriteFile(inputFile, []byte{}, 0644); err != nil {
				log.Printf("Failed to clear input file: %v", err)
			}

			email := parseEmail(string(data))
			entry := LogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Service:   "parser",
				Event:     "parsed_email",
				Data:      email,
			}

			jsonData, err := json.MarshalIndent(entry, "", "  ")
			if err != nil {
				log.Printf("Failed to encode JSON: %v", err)
				continue
			}

			log.Println(string(jsonData))
			fmt.Println(string(jsonData))

			if webhookURL != "" {
				if err := postWithRetry(webhookURL, jsonData, 5, 2*time.Second); err != nil {
					log.Printf("Failed to POST to webhook: %v", err)
				} else {
					log.Println("Posted to webhook successfully.")
				}
			}
			time.Sleep(intervalDur)
		}
	}

	log.Println("Parser service stopped gracefully.")
}
