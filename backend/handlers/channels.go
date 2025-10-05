package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"venotes/backend/database"
	"venotes/backend/models"

	"github.com/gin-gonic/gin"
)

// GetChannels returns all channels for a specific guild
func GetChannels(c *gin.Context) {
	guildID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guild ID"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.guild_id, c.name, c.created_at, c.updated_at 
		FROM channels c 
		WHERE c.guild_id = ? 
		ORDER BY c.name
	`, guildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(&channel.ID, &channel.GuildID, &channel.Name, &channel.CreatedAt, &channel.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		channels = append(channels, channel)
	}

	c.JSON(http.StatusOK, channels)
}

// GetChannel returns a specific channel by ID
func GetChannel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	var channel models.Channel
	err = database.DB.QueryRow(`
		SELECT id, guild_id, name, created_at, updated_at 
		FROM channels 
		WHERE id = ?
	`, id).Scan(&channel.ID, &channel.GuildID, &channel.Name, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, channel)
}

// CreateChannel creates a new channel in a guild
func CreateChannel(c *gin.Context) {
	guildID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guild ID"})
		return
	}

	var req models.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("INSERT INTO channels (guild_id, name) VALUES (?, ?)", guildID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var channel models.Channel
	err = database.DB.QueryRow(`
		SELECT id, guild_id, name, created_at, updated_at 
		FROM channels 
		WHERE id = ?
	`, id).Scan(&channel.ID, &channel.GuildID, &channel.Name, &channel.CreatedAt, &channel.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// DeleteChannel deletes a channel and all its logs
func DeleteChannel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM channels WHERE id = ?", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}
