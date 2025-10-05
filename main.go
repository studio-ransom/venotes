package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"venotes/backend/database"
	"venotes/backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var frontendAssets embed.FS

// getContentType returns the appropriate MIME type for a file based on its extension
func getContentType(filePath string) string {
	// Remove leading slash if present
	if strings.HasPrefix(filePath, "/") {
		filePath = filePath[1:]
	}
	
	ext := ""
	if dotIndex := strings.LastIndex(filePath, "."); dotIndex != -1 {
		ext = filePath[dotIndex:]
	}
	
	switch ext {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".ico":
		return "image/x-icon"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Initialize storage
	if err := handlers.InitStorage(); err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		// Guild routes
		api.GET("/guilds", handlers.GetGuilds)
		api.GET("/guilds/:id", handlers.GetGuild)
		api.POST("/guilds", handlers.CreateGuild)
		api.DELETE("/guilds/:id", handlers.DeleteGuild)

		// Channel routes - use different parameter names to avoid conflicts
		api.GET("/guilds/:id/channels", handlers.GetChannels)
		api.POST("/guilds/:id/channels", handlers.CreateChannel)
		api.GET("/channels/:id", handlers.GetChannel)
		api.DELETE("/channels/:id", handlers.DeleteChannel)

		// Log routes
		api.GET("/channels/:id/logs", handlers.GetLogs)
		api.POST("/channels/:id/logs", handlers.CreateLog)
		api.GET("/logs/:id", handlers.GetLog)
		api.PUT("/logs/:id", handlers.UpdateLog)
		api.DELETE("/logs/:id", handlers.DeleteLog)

	// File routes
	api.POST("/logs/:id/files", handlers.UploadFiles)
	api.GET("/files/:id", handlers.ServeFile)
	api.GET("/files/:id/content", handlers.GetFileContent)
	api.DELETE("/files/:id", handlers.DeleteFile)

	// Export/Import routes
	api.GET("/export", handlers.ExportData)
	api.POST("/import", handlers.ImportData)
	}

	// Create a sub-filesystem for the frontend assets
	frontendFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to create frontend filesystem:", err)
	}
	

	// Serve static files (frontend) from embedded assets with proper MIME types
	r.GET("/assets/*filepath", func(c *gin.Context) {
		filePath := c.Param("filepath")
		// Remove leading slash from filepath
		if strings.HasPrefix(filePath, "/") {
			filePath = filePath[1:]
		}
		
		// The embedded filesystem already has the assets/ prefix, so we need to add it
		fullPath := "assets/" + filePath
		
		file, err := frontendFS.Open(fullPath)
		if err != nil {
			log.Printf("File not found: %s, error: %v", fullPath, err)
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()
		
		// Set correct MIME type based on file extension
		contentType := getContentType(filePath)
		log.Printf("Serving file %s with content type %s", filePath, contentType)
		
		// Read file content into memory
		content, err := io.ReadAll(file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read file")
			return
		}
		
		c.Data(http.StatusOK, contentType, content)
	})
	
	// Serve index.html for root path
	r.GET("/", func(c *gin.Context) {
		indexFile, err := frontendFS.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Frontend not found")
			return
		}
		defer indexFile.Close()
		
		c.DataFromReader(http.StatusOK, -1, "text/html", indexFile, nil)
	})
	
	// Serve favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		faviconFile, err := frontendFS.Open("favicon.ico")
		if err != nil {
			c.String(http.StatusNotFound, "Favicon not found")
			return
		}
		defer faviconFile.Close()
		
		c.DataFromReader(http.StatusOK, -1, "image/x-icon", faviconFile, nil)
	})

	// Catch-all route for SPA
	r.NoRoute(func(c *gin.Context) {
		// Check if the request is for an API endpoint
		if c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		// Serve the SPA for all other routes from embedded assets
		indexFile, err := frontendFS.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Frontend not found")
			return
		}
		defer indexFile.Close()
		
		c.DataFromReader(http.StatusOK, -1, "text/html", indexFile, nil)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API available at http://localhost:%s/api", port)
	log.Printf("Frontend available at http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
