package handlers

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"venotes/backend/database"
	"venotes/backend/models"
	"venotes/backend/storage"
)

var fileStorage storage.Storage

// InitStorage initializes the storage backend
func InitStorage() error {
	factory := storage.NewStorageFactory()
	var err error
	fileStorage, err = factory.CreateStorage()
	return err
}

// UploadFiles handles file uploads for a log
func UploadFiles(c *gin.Context) {
	logIDStr := c.Param("id")
	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	// Verify log exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM logs WHERE id = ?)", logID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	var uploadedFiles []models.File

	for _, fileHeader := range files {
		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		// Calculate file hash
		hasher := sha256.New()
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			return
		}
		hasher.Write(fileContent)
		fileHash := fmt.Sprintf("%x", hasher.Sum(nil))

		// Check if file with this hash already exists in storage
		ext := filepath.Ext(fileHeader.Filename)
		deduplicatedFilename := fmt.Sprintf("%s%s", fileHash, ext)
		
		// Check if file already exists in storage
		exists, err := fileStorage.FileExists(c.Request.Context(), deduplicatedFilename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file existence"})
			return
		}
		
		// Only save file to storage if it doesn't already exist
		if !exists {
			// Create a reader from the file content
			contentReader := strings.NewReader(string(fileContent))
			
			// Store file using storage interface
			_, err = fileStorage.StoreFile(c.Request.Context(), deduplicatedFilename, contentReader)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
		}

		// Always create a new database entry (so it appears in logs)
		// Use a unique filename for the database entry but point to the deduplicated file
		uniqueFilename := fmt.Sprintf("%s_%d%s", fileHash[:16], time.Now().UnixNano(), ext)
		
		// Insert file record into database
		result, err := database.DB.Exec(`
			INSERT INTO files (log_id, filename, original_name, mime_type, size, path, hash)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, logID, uniqueFilename, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), 
			int64(len(fileContent)), deduplicatedFilename, fileHash)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file record"})
			return
		}

		fileID, _ := result.LastInsertId()
		uploadedFiles = append(uploadedFiles, models.File{
			ID:           int(fileID),
			LogID:        logID,
			Filename:     uniqueFilename,
			OriginalName: fileHeader.Filename,
			MimeType:     fileHeader.Header.Get("Content-Type"),
			Size:         int64(len(fileContent)),
			Path:         deduplicatedFilename, // Points to storage key
			Hash:         fileHash,
			CreatedAt:    time.Now(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{"files": uploadedFiles})
}

// ServeFile serves uploaded files
func ServeFile(c *gin.Context) {
	fileID := c.Param("id")
	
	var file models.File
	var hash sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, log_id, filename, original_name, mime_type, size, path, hash, created_at
		FROM files WHERE id = ?
	`, fileID).Scan(&file.ID, &file.LogID, &file.Filename, &file.OriginalName, &file.MimeType, &file.Size, &file.Path, &hash, &file.CreatedAt)
	
	if err == nil {
		if hash.Valid {
			file.Hash = hash.String
		} else {
			file.Hash = ""
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if file exists in storage
	exists, err := fileStorage.FileExists(c.Request.Context(), file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check file existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in storage"})
		return
	}

	// Get file from storage
	fileReader, err := fileStorage.GetFile(c.Request.Context(), file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer fileReader.Close()

	// Set appropriate headers
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.OriginalName))
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// Copy file content to response
	_, err = io.Copy(c.Writer, fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serve file"})
		return
	}
}

// GetFileContent returns the content of a text file
func GetFileContent(c *gin.Context) {
	fileID := c.Param("id")
	
	var file models.File
	var hash sql.NullString
	err := database.DB.QueryRow(`
		SELECT id, log_id, filename, original_name, mime_type, size, path, hash, created_at
		FROM files WHERE id = ?
	`, fileID).Scan(&file.ID, &file.LogID, &file.Filename, &file.OriginalName, &file.MimeType, &file.Size, &file.Path, &hash, &file.CreatedAt)
	
	if err == nil {
		if hash.Valid {
			file.Hash = hash.String
		} else {
			file.Hash = ""
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if it's a text file
	if !isTextFile(file.MimeType, file.OriginalName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is not a text file"})
		return
	}

	// Get file from storage
	fileReader, err := fileStorage.GetFile(c.Request.Context(), file.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer fileReader.Close()

	// Read file content
	content, err := io.ReadAll(fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, string(content))
}

// DeleteFile deletes a file
func DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	
	var file models.File
	err := database.DB.QueryRow(`
		SELECT id, log_id, filename, original_name, mime_type, size, path, hash, created_at
		FROM files WHERE id = ?
	`, fileID).Scan(&file.ID, &file.LogID, &file.Filename, &file.OriginalName, &file.MimeType, &file.Size, &file.Path, &file.Hash, &file.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Delete from database
	_, err = database.DB.Exec("DELETE FROM files WHERE id = ?", fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	// Delete from storage
	if err := fileStorage.DeleteFile(c.Request.Context(), file.Path); err != nil {
		// Log error but don't fail the request since DB record is deleted
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// Helper function to generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// Helper function to check if file is text
func isTextFile(mimeType, filename string) bool {
	textMimeTypes := []string{
		"text/",
		"application/json",
		"application/javascript",
		"application/xml",
	}
	
	textExtensions := []string{
		".txt", ".md", ".json", ".js", ".css", ".html", ".xml", ".yaml", ".yml",
		".py", ".go", ".java", ".cpp", ".c", ".h", ".php", ".rb", ".rs",
	}
	
	// Check MIME type
	for _, mime := range textMimeTypes {
		if strings.HasPrefix(mimeType, mime) {
			return true
		}
	}
	
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	
	return false
}
