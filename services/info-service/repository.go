package main

import (
	"database/sql"
	"time"

	"github.com/son/bot-cafe/shared"
)

// Repository handles database operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// InitSchema initializes database schema
func (r *Repository) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS cafe_info (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		address TEXT NOT NULL,
		phone TEXT NOT NULL,
		email TEXT,
		opening_hour TEXT NOT NULL,
		closing_hour TEXT NOT NULL,
		description TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Insert default café info if not exists
	INSERT OR IGNORE INTO cafe_info (id, name, address, phone, opening_hour, closing_hour, description)
	VALUES (1, 'Café/Resto Bot', 'Jl. Contoh No. 123', '081234567890', '08:00', '22:00', 'Selamat datang di Café kami!');
	`
	return shared.ExecuteSchema(r.db, schema)
}

// GetCafeInfo gets café information
func (r *Repository) GetCafeInfo() (*CafeInfo, error) {
	query := `SELECT id, name, address, phone, email, opening_hour, closing_hour, description, updated_at 
			  FROM cafe_info WHERE id = 1`
	var info CafeInfo
	err := r.db.QueryRow(query).Scan(
		&info.ID, &info.Name, &info.Address, &info.Phone, &info.Email,
		&info.OpeningHour, &info.ClosingHour, &info.Description, &info.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Informasi café")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &info, nil
}

// UpdateCafeInfo updates café information

func (r *Repository) UpdateCafeInfo(info *CafeInfo) error {
	query := `UPDATE cafe_info SET name = ?, address = ?, phone = ?, email = ?, 
			  opening_hour = ?, closing_hour = ?, description = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = 1`
	_, err := r.db.Exec(query, info.Name, info.Address, info.Phone, info.Email,
		info.OpeningHour, info.ClosingHour, info.Description)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	info.UpdatedAt = time.Now()
	return nil
}
