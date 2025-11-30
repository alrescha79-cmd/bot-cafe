package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/son/bot-cafe/shared"
)

func main() {
	// Load environment variables
	godotenv.Load()

	port := os.Getenv("PROMO_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	dbPath := os.Getenv("PROMO_DB_PATH")
	if dbPath == "" {
		dbPath = "./data/promo.db"
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
	shared.LogInfo("Promo service starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
