package models

import "time"

// Guild represents a workspace/category
type Guild struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Channel represents a topic within a guild
type Channel struct {
	ID        int       `json:"id" db:"id"`
	GuildID   int       `json:"guild_id" db:"guild_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Log represents a note entry within a channel
type Log struct {
	ID        int       `json:"id" db:"id"`
	ChannelID int       `json:"channel_id" db:"channel_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Files     []File    `json:"files,omitempty"`
}

// File represents an uploaded file attachment
type File struct {
	ID        int       `json:"id" db:"id"`
	LogID     int       `json:"log_id" db:"log_id"`
	Filename  string    `json:"filename" db:"filename"`
	OriginalName string `json:"original_name" db:"original_name"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Size      int64     `json:"size" db:"size"`
	Path      string    `json:"path" db:"path"`
	Hash      string    `json:"hash" db:"hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateGuildRequest represents the request to create a guild
type CreateGuildRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateChannelRequest represents the request to create a channel
type CreateChannelRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateLogRequest represents the request to create a log
type CreateLogRequest struct {
	Content string `json:"content"`
}

// UpdateLogRequest represents the request to update a log
type UpdateLogRequest struct {
	Content string `json:"content" binding:"required"`
}
