package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.SetOutput(os.Stdout)
	log.Println("Webhook service starting on port", port)

	// Handler for any of these routes
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST supported", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			log.Println("Failed to read body:", err)
			return
		}

		log.Printf("[%s] Received webhook on %s:\n%s\n", time.Now().Format(time.RFC3339), r.URL.Path, string(body))
		w.WriteHeader(http.StatusOK)
	}

	// Register desired routes (no catch-all "/" to avoid matching unrelated requests)
	http.HandleFunc("/webhook", handler)
	http.HandleFunc("/event", handler)
	http.HandleFunc("/receive", handler)
	http.HandleFunc("/email", handler)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
