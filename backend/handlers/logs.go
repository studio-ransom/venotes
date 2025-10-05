package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"venotes/backend/database"
	"venotes/backend/models"

	"github.com/gin-gonic/gin"
)

// GetLogs returns all logs for a specific channel
func GetLogs(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT l.id, l.channel_id, l.content, l.created_at, l.updated_at 
		FROM logs l 
		WHERE l.channel_id = ? 
		ORDER BY l.created_at DESC
	`, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(&log.ID, &log.ChannelID, &log.Content, &log.CreatedAt, &log.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Load files for this log
		files, err := getFilesForLog(log.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Files = files
		
		logs = append(logs, log)
	}

	c.JSON(http.StatusOK, logs)
}

// GetLog returns a specific log by ID
func GetLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	var log models.Log
	err = database.DB.QueryRow(`
		SELECT id, channel_id, content, created_at, updated_at 
		FROM logs 
		WHERE id = ?
	`, id).Scan(&log.ID, &log.ChannelID, &log.Content, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, log)
}

// CreateLog creates a new log in a channel
func CreateLog(c *gin.Context) {
	channelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var req models.CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("INSERT INTO logs (channel_id, content) VALUES (?, ?)", channelID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load the created log with files
	var log models.Log
	err = database.DB.QueryRow(`
		SELECT id, channel_id, content, created_at, updated_at 
		FROM logs 
		WHERE id = ?
	`, id).Scan(&log.ID, &log.ChannelID, &log.Content, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load files for the log
	files, err := getFilesForLog(log.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Files = files

	c.JSON(http.StatusCreated, log)
}

// UpdateLog updates an existing log
func UpdateLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	var req models.UpdateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("UPDATE logs SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	var log models.Log
	err = database.DB.QueryRow(`
		SELECT id, channel_id, content, created_at, updated_at 
		FROM logs 
		WHERE id = ?
	`, id).Scan(&log.ID, &log.ChannelID, &log.Content, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// DeleteLog deletes a log
func DeleteLog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM logs WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Log deleted successfully"})
}

// getFilesForLog retrieves all files associated with a log
func getFilesForLog(logID int) ([]models.File, error) {
	rows, err := database.DB.Query(`
		SELECT id, log_id, filename, original_name, mime_type, size, path, hash, created_at
		FROM files WHERE log_id = ?
		ORDER BY created_at ASC
	`, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var hash sql.NullString
		err := rows.Scan(&file.ID, &file.LogID, &file.Filename, &file.OriginalName, &file.MimeType, &file.Size, &file.Path, &hash, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		if hash.Valid {
			file.Hash = hash.String
		} else {
			file.Hash = ""
		}
		files = append(files, file)
	}

	return files, nil
}
