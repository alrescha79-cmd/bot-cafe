package main

import "time"

// Admin represents an admin user
type Admin struct {
	ID         int       `json:"id"`
	TelegramID string    `json:"telegram_id"`
	Username   string    `json:"username"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// Session represents an admin session
type Session struct {
	ID        int       `json:"id"`
	AdminID   int       `json:"admin_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
