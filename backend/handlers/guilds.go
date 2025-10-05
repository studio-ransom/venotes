package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"venotes/backend/database"
	"venotes/backend/models"

	"github.com/gin-gonic/gin"
)

// GetGuilds returns all guilds
func GetGuilds(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, created_at, updated_at FROM guilds ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var guilds []models.Guild
	for rows.Next() {
		var guild models.Guild
		err := rows.Scan(&guild.ID, &guild.Name, &guild.CreatedAt, &guild.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		guilds = append(guilds, guild)
	}

	c.JSON(http.StatusOK, guilds)
}

// GetGuild returns a specific guild by ID
func GetGuild(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guild ID"})
		return
	}

	var guild models.Guild
	err = database.DB.QueryRow("SELECT id, name, created_at, updated_at FROM guilds WHERE id = ?", id).Scan(
		&guild.ID, &guild.Name, &guild.CreatedAt, &guild.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guild not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, guild)
}

// CreateGuild creates a new guild
func CreateGuild(c *gin.Context) {
	var req models.CreateGuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("INSERT INTO guilds (name) VALUES (?)", req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var guild models.Guild
	err = database.DB.QueryRow("SELECT id, name, created_at, updated_at FROM guilds WHERE id = ?", id).Scan(
		&guild.ID, &guild.Name, &guild.CreatedAt, &guild.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, guild)
}

// DeleteGuild deletes a guild and all its channels and logs
func DeleteGuild(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guild ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM guilds WHERE id = ?", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Guild not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guild deleted successfully"})
}
