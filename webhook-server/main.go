package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Email struct {
	Subject string `json:"subject"`
	From    string `json:"from"`
	To      string `json:"to"`
	Date    string `json:"date"`
	Body    string `json:"body"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Event     string `json:"event"`
	Data      Email  `json:"data"`
}

func writeLog(entry LogEntry, logPath string) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("error marshaling log entry: %v", err)
		return
	}
	file.WriteString(string(jsonData) + "\n")
}

func emailHandler(logPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "invalid method", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}

		var email Email
		if err := json.Unmarshal(body, &email); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		entry := LogEntry{
			Timestamp: time.Now().Format("2006-01-02 15:04:05 MST"),
			Service:   "webhook-server",
			Event:     "email_received",
			Data:      email,
		}
		writeLog(entry, logPath)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Email received and logged")
	}
}

func main() {
	// Create logs/ directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	logPath := os.Getenv("LOG_FILE_PATH")
	if logPath == "" {
		logPath = "logs/webhook-server.log"
		log.Printf("LOG_FILE_PATH not set, defaulting to %s", logPath)
	}

	http.HandleFunc("/email", emailHandler(logPath))

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	fmt.Printf("Webhook server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
