package shared

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database
func InitDB(dbPath string) (*sql.DB, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database initialized: %s", dbPath)
	return db, nil
}

// ExecuteSchema executes SQL schema
func ExecuteSchema(db *sql.DB, schema string) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
