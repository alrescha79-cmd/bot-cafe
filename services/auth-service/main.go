package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/son/bot-cafe/shared"
)

func main() {
	// Load environment variables
	godotenv.Load()

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	dbPath := os.Getenv("AUTH_DB_PATH")
	if dbPath == "" {
		dbPath = "./data/auth.db"
	}

	// Initialize database
	db, err := shared.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := NewRepository(db)
	if err := repo.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Cleanup expired sessions periodically
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := repo.CleanupExpiredSessions(); err != nil {
				shared.LogError("Failed to cleanup sessions: %v", err)
			}
		}
	}()

	// Initialize handler
	handler := NewHandler(repo)

	// Setup routes
	http.HandleFunc("/", handler.HandleRequest)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", port)
	shared.LogInfo("Auth service starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
