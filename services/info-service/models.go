package main

import "time"

// CafeInfo represents café information
type CafeInfo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	OpeningHour string    `json:"opening_hour"`
	ClosingHour string    `json:"closing_hour"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}
