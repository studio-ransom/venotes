package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	dbPath := filepath.Join(dataDir, "notes.db")
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Insert default data
	if err := insertDefaultData(); err != nil {
		log.Printf("Warning: failed to insert default data: %v", err)
	}

	log.Printf("Database initialized successfully at %s", dbPath)
	return nil
}

// createTables creates the necessary tables
func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS guilds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			guild_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (guild_id) REFERENCES guilds (id) ON DELETE CASCADE,
			UNIQUE(guild_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			channel_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_id) REFERENCES channels (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			log_id INTEGER NOT NULL,
			filename TEXT NOT NULL,
			original_name TEXT NOT NULL,
			mime_type TEXT NOT NULL,
			size INTEGER NOT NULL,
			path TEXT NOT NULL,
			hash TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (log_id) REFERENCES logs (id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_channels_guild_id ON channels(guild_id)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_channel_id ON logs(channel_id)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_files_log_id ON files(log_id)`,
		`CREATE INDEX IF NOT EXISTS idx_files_hash ON files(hash)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	// Add hash column to existing files table if it doesn't exist
	var columnExists int
	err := DB.QueryRow("SELECT COUNT(*) FROM pragma_table_info('files') WHERE name='hash'").Scan(&columnExists)
	if err == nil && columnExists == 0 {
		// Hash column doesn't exist, add it
		_, err = DB.Exec("ALTER TABLE files ADD COLUMN hash TEXT")
		if err != nil {
			// Log error but don't fail initialization
			fmt.Printf("Warning: Could not add hash column: %v\n", err)
		}
	}

	return nil
}

// insertDefaultData inserts default guilds and channels
func insertDefaultData() error {
	// Check if data already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM guilds").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Data already exists
	}

	// Insert default guilds and channels
	defaultData := []struct {
		guildName    string
		channelNames []string
	}{
		{"Work", []string{"Coding", "Ideas", "Meetings"}},
		{"Personal", []string{"Thoughts", "Ideas", "Journal"}},
	}

	for _, data := range defaultData {
		// Insert guild
		result, err := DB.Exec("INSERT INTO guilds (name) VALUES (?)", data.guildName)
		if err != nil {
			return err
		}

		guildID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert channels for this guild
		for _, channelName := range data.channelNames {
			_, err := DB.Exec("INSERT INTO channels (guild_id, name) VALUES (?, ?)", guildID, channelName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
