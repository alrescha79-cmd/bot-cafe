package main

import "time"

// Promo represents a promotional offer
type Promo struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Discount     int       `json:"discount"`      // percentage or amount
	DiscountType string    `json:"discount_type"` // "percentage" or "amount"
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
