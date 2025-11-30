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
	CREATE TABLE IF NOT EXISTS promos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		discount INTEGER NOT NULL,
		discount_type TEXT NOT NULL,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_active ON promos(is_active);
	CREATE INDEX IF NOT EXISTS idx_dates ON promos(start_date, end_date);
	`
	return shared.ExecuteSchema(r.db, schema)
}

// CreatePromo creates a new promo
func (r *Repository) CreatePromo(promo *Promo) (*Promo, error) {
	query := `INSERT INTO promos (title, description, discount, discount_type, start_date, end_date, is_active) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, promo.Title, promo.Description, promo.Discount, promo.DiscountType,
		promo.StartDate, promo.EndDate, promo.IsActive)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	id, _ := result.LastInsertId()
	promo.ID = int(id)
	promo.CreatedAt = time.Now()
	promo.UpdatedAt = time.Now()
	return promo, nil
}

// GetPromoByID gets promo by ID
func (r *Repository) GetPromoByID(id int) (*Promo, error) {
	query := `SELECT id, title, description, discount, discount_type, start_date, end_date, is_active, created_at, updated_at 
			  FROM promos WHERE id = ?`
	var promo Promo
	err := r.db.QueryRow(query, id).Scan(
		&promo.ID, &promo.Title, &promo.Description, &promo.Discount, &promo.DiscountType,
		&promo.StartDate, &promo.EndDate, &promo.IsActive, &promo.CreatedAt, &promo.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Promo")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &promo, nil
}

// ListPromos lists all promos with optional filters
func (r *Repository) ListPromos(activeOnly bool) ([]Promo, error) {
	query := `SELECT id, title, description, discount, discount_type, start_date, end_date, is_active, created_at, updated_at 
			  FROM promos WHERE 1=1`
	if activeOnly {
		query += ` AND is_active = 1 AND start_date <= datetime('now') AND end_date >= datetime('now')`
	}
	query += ` ORDER BY start_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	defer rows.Close()

	var promos []Promo
	for rows.Next() {
		var promo Promo
		if err := rows.Scan(&promo.ID, &promo.Title, &promo.Description, &promo.Discount, &promo.DiscountType,
			&promo.StartDate, &promo.EndDate, &promo.IsActive, &promo.CreatedAt, &promo.UpdatedAt); err != nil {
			return nil, shared.NewDatabaseError(err)
		}
		promos = append(promos, promo)
	}
	return promos, nil
}

// UpdatePromo updates a promo
func (r *Repository) UpdatePromo(promo *Promo) error {
	query := `UPDATE promos SET title = ?, description = ?, discount = ?, discount_type = ?, 
			  start_date = ?, end_date = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	result, err := r.db.Exec(query, promo.Title, promo.Description, promo.Discount, promo.DiscountType,
		promo.StartDate, promo.EndDate, promo.IsActive, promo.ID)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Promo")
	}
	return nil
}

// DeletePromo deletes a promo
func (r *Repository) DeletePromo(id int) error {
	query := `DELETE FROM promos WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Promo")
	}
	return nil
}
