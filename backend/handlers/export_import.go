package handlers

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"venotes/backend/database"
)

// ExportData exports all data to a zip file
func ExportData(c *gin.Context) {
	password := c.Query("password")
	encrypt := password != ""

	// Create a temporary zip file
	zipPath := fmt.Sprintf("data/export_%d.zip", time.Now().Unix())
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create export file"})
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Export database data
	if err := exportDatabaseData(zipWriter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export database: %v", err)})
		return
	}

	// Export uploaded files
	if err := exportUploadedFiles(zipWriter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export files: %v", err)})
		return
	}

	zipWriter.Close()
	zipFile.Close()

	// If encryption is requested, encrypt the zip file
	if encrypt {
		encryptedPath := zipPath + ".enc"
		if err := encryptFile(zipPath, encryptedPath, password); err != nil {
			os.Remove(zipPath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to encrypt file: %v", err)})
			return
		}
		os.Remove(zipPath)
		zipPath = encryptedPath
	}

	// Serve the file
	filename := fmt.Sprintf("venotes-export-%s", time.Now().Format("2006-01-02"))
	if encrypt {
		filename += ".enc"
		c.Header("Content-Type", "application/octet-stream")
	} else {
		filename += ".zip"
		c.Header("Content-Type", "application/zip")
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.File(zipPath)

	// Clean up the temporary file
	os.Remove(zipPath)
}

// ImportData imports data from a zip file
func ImportData(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	password := c.PostForm("password")
	encrypted := password != ""

	// Save uploaded file temporarily
	tempPath := fmt.Sprintf("data/import_%d", time.Now().Unix())
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}
	defer os.Remove(tempPath)

	// If encrypted, decrypt first
	if encrypted {
		decryptedPath := tempPath + "_decrypted"
		if err := decryptFile(tempPath, decryptedPath, password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password or corrupted file"})
			return
		}
		tempPath = decryptedPath
		defer os.Remove(tempPath)
	}

	// Import the data
	if err := importDataFromZip(tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to import data: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data imported successfully"})
}

// exportDatabaseData exports all database tables to JSON files
func exportDatabaseData(zipWriter *zip.Writer) error {
	// Export guilds
	if err := exportTable(zipWriter, "guilds.json", "SELECT * FROM guilds"); err != nil {
		return err
	}

	// Export channels
	if err := exportTable(zipWriter, "channels.json", "SELECT * FROM channels"); err != nil {
		return err
	}

	// Export logs
	if err := exportTable(zipWriter, "logs.json", "SELECT * FROM logs"); err != nil {
		return err
	}

	// Export files
	if err := exportTable(zipWriter, "files.json", "SELECT * FROM files"); err != nil {
		return err
	}

	return nil
}

// exportTable exports a table to a JSON file in the zip
func exportTable(zipWriter *zip.Writer, filename, query string) error {
	rows, err := database.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	// Create JSON file in zip
	file, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	return err
}

// exportUploadedFiles exports all uploaded files to the zip
func exportUploadedFiles(zipWriter *zip.Writer) error {
	rows, err := database.DB.Query("SELECT path, original_name FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path, originalName string
		if err := rows.Scan(&path, &originalName); err != nil {
			continue
		}

		// Read the file
		fileData, err := os.ReadFile(path)
		if err != nil {
			continue // Skip files that don't exist
		}

		// Create file in zip with original name
		zipFile, err := zipWriter.Create(fmt.Sprintf("uploads/%s", originalName))
		if err != nil {
			continue
		}

		zipFile.Write(fileData)
	}

	return nil
}

// importDataFromZip imports data from a zip file
func importDataFromZip(zipPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Import database tables
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, ".json") {
			if err := importTable(file); err != nil {
				return err
			}
		} else if strings.HasPrefix(file.Name, "uploads/") {
			if err := importUploadedFile(file); err != nil {
				return err
			}
		}
	}

	return nil
}

// importTable imports a JSON table file
func importTable(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}

	tableName := strings.TrimSuffix(file.Name, ".json")
	
	// Import based on table name
	switch tableName {
	case "guilds":
		return importGuilds(records)
	case "channels":
		return importChannels(records)
	case "logs":
		return importLogs(records)
	case "files":
		return importFiles(records)
	}

	return nil
}

// importGuilds imports guild records
func importGuilds(records []map[string]interface{}) error {
	for _, record := range records {
		name, ok := record["name"].(string)
		if !ok {
			continue
		}

		// Check if guild already exists
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM guilds WHERE name = ?", name).Scan(&count)
		if err != nil || count > 0 {
			continue // Skip if exists or error
		}

		// Insert new guild
		database.DB.Exec("INSERT INTO guilds (name) VALUES (?)", name)
	}
	return nil
}

// importChannels imports channel records
func importChannels(records []map[string]interface{}) error {
	for _, record := range records {
		name, ok := record["name"].(string)
		if !ok {
			continue
		}

		guildID, ok := record["guild_id"].(float64)
		if !ok {
			continue
		}

		// Check if channel already exists
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM channels WHERE name = ? AND guild_id = ?", name, int(guildID)).Scan(&count)
		if err != nil || count > 0 {
			continue // Skip if exists or error
		}

		// Insert new channel
		database.DB.Exec("INSERT INTO channels (name, guild_id) VALUES (?, ?)", name, int(guildID))
	}
	return nil
}

// importLogs imports log records
func importLogs(records []map[string]interface{}) error {
	for _, record := range records {
		content, ok := record["content"].(string)
		if !ok {
			continue
		}

		channelID, ok := record["channel_id"].(float64)
		if !ok {
			continue
		}

		// Check if log with same content and channel already exists
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM logs WHERE content = ? AND channel_id = ?", content, int(channelID)).Scan(&count)
		if err != nil || count > 0 {
			continue // Skip if exists or error
		}

		// Insert new log
		database.DB.Exec("INSERT INTO logs (channel_id, content) VALUES (?, ?)", int(channelID), content)
	}
	return nil
}

// importFiles imports file records
func importFiles(records []map[string]interface{}) error {
	for _, record := range records {
		originalName, ok := record["original_name"].(string)
		if !ok {
			continue
		}

		hash, ok := record["hash"].(string)
		if !ok {
			continue
		}

		// Check if file with this hash already exists
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM files WHERE hash = ?", hash).Scan(&count)
		if err != nil || count > 0 {
			continue // Skip if exists or error
		}

		// Get the log ID from the record (if available) or find a matching log
		var logID int
		if logIDFloat, ok := record["log_id"].(float64); ok {
			logID = int(logIDFloat)
		} else {
			// If no log_id in record, try to find a matching log by content and channel
			channelID, ok := record["channel_id"].(float64)
			if !ok {
				continue
			}
			
			// Try to find a log that matches the file's original context
			// For now, just get the latest log in the channel
			err = database.DB.QueryRow("SELECT id FROM logs WHERE channel_id = ? ORDER BY created_at DESC LIMIT 1", int(channelID)).Scan(&logID)
			if err != nil {
				continue
			}
		}

		// Insert file record
		database.DB.Exec(`
			INSERT INTO files (log_id, filename, original_name, mime_type, size, path, hash)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, logID, record["filename"], originalName, record["mime_type"], record["size"], record["path"], hash)
	}
	return nil
}

// importUploadedFile imports an uploaded file
func importUploadedFile(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create uploads directory if it doesn't exist
	os.MkdirAll("data/uploads", 0755)

	// Extract file to uploads directory
	path := filepath.Join("data/uploads", filepath.Base(file.Name))
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, rc)
	return err
}

// encryptFile encrypts a file with AES-256-GCM
func encryptFile(inputPath, outputPath, password string) error {
	// Read the file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Create key from password
	key := sha256.Sum256([]byte(password))

	// Create cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Write encrypted file
	return os.WriteFile(outputPath, ciphertext, 0644)
}

// decryptFile decrypts a file with AES-256-GCM
func decryptFile(inputPath, outputPath, password string) error {
	// Read encrypted file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Create key from password
	key := sha256.Sum256([]byte(password))

	// Create cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	// Write decrypted file
	return os.WriteFile(outputPath, plaintext, 0644)
}
