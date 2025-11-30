package main

import (
	"database/sql"
	"time"

	"github.com/alrescha79-cmd/bot-cafe/shared"
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
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS menus (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		price INTEGER NOT NULL,
		category TEXT NOT NULL,
		photo_url TEXT,
		is_available BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category) REFERENCES categories(name)
	);

	CREATE INDEX IF NOT EXISTS idx_category ON menus(category);
	CREATE INDEX IF NOT EXISTS idx_available ON menus(is_available);

	-- Insert default categories
	INSERT OR IGNORE INTO categories (name) VALUES ('Makanan');
	INSERT OR IGNORE INTO categories (name) VALUES ('Minuman');
	INSERT OR IGNORE INTO categories (name) VALUES ('Snack');
	INSERT OR IGNORE INTO categories (name) VALUES ('Coffee');
	`
	return shared.ExecuteSchema(r.db, schema)
}

// CreateMenu creates a new menu
func (r *Repository) CreateMenu(menu *Menu) (*Menu, error) {
	query := `INSERT INTO menus (name, description, price, category, photo_url, is_available) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, menu.Name, menu.Description, menu.Price, menu.Category, menu.PhotoURL, menu.IsAvailable)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	id, _ := result.LastInsertId()
	menu.ID = int(id)
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()
	return menu, nil
}

// GetMenuByID gets menu by ID
func (r *Repository) GetMenuByID(id int) (*Menu, error) {
	query := `SELECT id, name, description, price, category, photo_url, is_available, created_at, updated_at 
			  FROM menus WHERE id = ?`
	var menu Menu
	err := r.db.QueryRow(query, id).Scan(
		&menu.ID, &menu.Name, &menu.Description, &menu.Price, &menu.Category,
		&menu.PhotoURL, &menu.IsAvailable, &menu.CreatedAt, &menu.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Menu")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &menu, nil
}

// ListMenus lists all menus with optional filters
func (r *Repository) ListMenus(category string, availableOnly bool) ([]Menu, error) {
	query := `SELECT id, name, description, price, category, photo_url, is_available, created_at, updated_at 
			  FROM menus WHERE 1=1`
	args := []interface{}{}

	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}

	if availableOnly {
		query += ` AND is_available = 1`
	}

	query += ` ORDER BY category, name`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price, &menu.Category,
			&menu.PhotoURL, &menu.IsAvailable, &menu.CreatedAt, &menu.UpdatedAt); err != nil {
			return nil, shared.NewDatabaseError(err)
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

// UpdateMenu updates a menu
func (r *Repository) UpdateMenu(menu *Menu) error {
	query := `UPDATE menus SET name = ?, description = ?, price = ?, category = ?, 
			  photo_url = ?, is_available = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	result, err := r.db.Exec(query, menu.Name, menu.Description, menu.Price, menu.Category,
		menu.PhotoURL, menu.IsAvailable, menu.ID)
	if err != nil {
		return shared.NewDatabaseError(err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Menu")
	}
	return nil
}

// DeleteMenu deletes a menu
func (r *Repository) DeleteMenu(id int) error {
	query := `DELETE FROM menus WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError(err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Menu")
	}
	return nil
}

// ListCategories lists all categories
func (r *Repository) ListCategories() ([]Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			return nil, shared.NewDatabaseError(err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// CreateCategory creates a new category
func (r *Repository) CreateCategory(name string) (*Category, error) {
	query := `INSERT INTO categories (name) VALUES (?)`
	result, err := r.db.Exec(query, name)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	id, _ := result.LastInsertId()
	return &Category{
		ID:        int(id),
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}

// DeleteCategory deletes a category
func (r *Repository) DeleteCategory(name string) error {
	// Check if category has menus
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM menus WHERE category = ?`, name).Scan(&count)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	if count > 0 {
		return shared.NewInvalidInputError("Kategori masih memiliki menu, tidak dapat dihapus")
	}

	query := `DELETE FROM categories WHERE name = ?`
	result, err := r.db.Exec(query, name)
	if err != nil {
		return shared.NewDatabaseError(err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Kategori")
	}
	return nil
}
