package main

import "time"

// Media represents a media file
type Media struct {
	ID         int       `json:"id"`
	FileName   string    `json:"file_name"`
	FileURL    string    `json:"file_url"`
	FileType   string    `json:"file_type"`
	EntityID   int       `json:"entity_id"`   // ID of menu/promo
	EntityType string    `json:"entity_type"` // "menu" or "promo"
	CreatedAt  time.Time `json:"created_at"`
}
