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
	CREATE TABLE IF NOT EXISTS media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_name TEXT NOT NULL,
		file_url TEXT NOT NULL,
		file_type TEXT NOT NULL,
		entity_id INTEGER,
		entity_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_entity ON media(entity_id, entity_type);
	`
	return shared.ExecuteSchema(r.db, schema)
}

// CreateMedia creates a new media record
func (r *Repository) CreateMedia(media *Media) (*Media, error) {
	query := `INSERT INTO media (file_name, file_url, file_type, entity_id, entity_type) 
			  VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, media.FileName, media.FileURL, media.FileType, media.EntityID, media.EntityType)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}

	id, _ := result.LastInsertId()
	media.ID = int(id)
	media.CreatedAt = time.Now()
	return media, nil
}

// GetMediaByID gets media by ID
func (r *Repository) GetMediaByID(id int) (*Media, error) {
	query := `SELECT id, file_name, file_url, file_type, entity_id, entity_type, created_at 
			  FROM media WHERE id = ?`
	var media Media
	err := r.db.QueryRow(query, id).Scan(
		&media.ID, &media.FileName, &media.FileURL, &media.FileType,
		&media.EntityID, &media.EntityType, &media.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shared.NewNotFoundError("Media")
	}
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	return &media, nil
}

// ListMediaByEntity lists media by entity
func (r *Repository) ListMediaByEntity(entityID int, entityType string) ([]Media, error) {
	query := `SELECT id, file_name, file_url, file_type, entity_id, entity_type, created_at 
			  FROM media WHERE entity_id = ? AND entity_type = ? ORDER BY created_at DESC`
	rows, err := r.db.Query(query, entityID, entityType)
	if err != nil {
		return nil, shared.NewDatabaseError(err)
	}
	defer rows.Close()

	var medias []Media
	for rows.Next() {
		var media Media
		if err := rows.Scan(&media.ID, &media.FileName, &media.FileURL, &media.FileType,
			&media.EntityID, &media.EntityType, &media.CreatedAt); err != nil {
			return nil, shared.NewDatabaseError(err)
		}
		medias = append(medias, media)
	}
	return medias, nil
}

// DeleteMedia deletes a media
func (r *Repository) DeleteMedia(id int) error {
	query := `DELETE FROM media WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError(err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return shared.NewNotFoundError("Media")
	}
	return nil
}
