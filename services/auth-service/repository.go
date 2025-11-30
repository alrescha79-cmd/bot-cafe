package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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
	CREATE TABLE IF NOT EXISTS admins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id TEXT UNIQUE NOT NULL,
		username TEXT NOT NULL,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		admin_id INTEGER NOT NULL,
		token TEXT UNIQUE NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (admin_id) REFERENCES admins(id)
	);

	CREATE INDEX IF NOT EXISTS idx_telegram_id ON admins(telegram_id);
	CREATE INDEX IF NOT EXISTS idx_token ON sessions(token);
	`
	return shared.ExecuteSchema(r.db, schema)
}

// CreateAdmin creates a new admin
func (r *Repository) CreateAdmin(telegramID, username string) (*Admin, error) {
	query := `INSERT INTO admins (telegram_id, username) VALUES (?, ?)`
	result, err := r.db.Exec(query, telegramID, username)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	id, _ := result.LastInsertId()
	return &Admin{
		ID:         int(id),
		TelegramID: telegramID,
		Username:   username,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}, nil
}

// GetAdminByTelegramID gets admin by telegram ID
func (r *Repository) GetAdminByTelegramID(telegramID string) (*Admin, error) {
	query := `SELECT id, telegram_id, username, is_active, created_at FROM admins WHERE telegram_id = ?`
	var admin Admin
	err := r.db.QueryRow(query, telegramID).Scan(
		&admin.ID, &admin.TelegramID, &admin.Username, &admin.IsActive, &admin.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Admin")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &admin, nil
}

// GetAdminByUsername gets admin by username
func (r *Repository) GetAdminByUsername(username string) (*Admin, error) {
	query := `SELECT id, telegram_id, username, is_active, created_at FROM admins WHERE username = ?`
	var admin Admin
	err := r.db.QueryRow(query, username).Scan(
		&admin.ID, &admin.TelegramID, &admin.Username, &admin.IsActive, &admin.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Admin")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &admin, nil
}

// ListAdmins lists all admins
func (r *Repository) ListAdmins() ([]Admin, error) {
	query := `SELECT id, telegram_id, username, is_active, created_at FROM admins ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.ID, &admin.TelegramID, &admin.Username, &admin.IsActive, &admin.CreatedAt); err != nil {
			return nil, shared.NewDatabaseError(err)
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

// CreateSession creates a new session
func (r *Repository) CreateSession(adminID int) (*Session, error) {
	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	query := `INSERT INTO sessions (admin_id, token, expires_at) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, adminID, token, expiresAt)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	id, _ := result.LastInsertId()
	return &Session{
		ID:        int(id),
		AdminID:   adminID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

// VerifySession verifies a session token
func (r *Repository) VerifySession(token string) (*Admin, error) {
	query := `
		SELECT a.id, a.telegram_id, a.username, a.is_active, a.created_at 
		FROM admins a
		JOIN sessions s ON a.id = s.admin_id
		WHERE s.token = ? AND s.expires_at > datetime('now') AND a.is_active = 1
	`
	var admin Admin
	err := r.db.QueryRow(query, token).Scan(
		&admin.ID, &admin.TelegramID, &admin.Username, &admin.IsActive, &admin.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewUnauthorizedError()
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &admin, nil
}

// DeleteSession deletes a session
func (r *Repository) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := r.db.Exec(query, token)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions
func (r *Repository) CleanupExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < datetime('now')`
	_, err := r.db.Exec(query)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	return nil
}

// UpdateAdminStatus updates admin active status
func (r *Repository) UpdateAdminStatus(telegramID string, isActive bool) error {
	query := `UPDATE admins SET is_active = ? WHERE telegram_id = ?`
	_, err := r.db.Exec(query, isActive, telegramID)
	if err != nil {
		return shared.NewDatabaseError(err)
	}
	return nil
}

// generateToken generates a random token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
