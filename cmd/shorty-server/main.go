package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/database"
)

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("SHORTY_DB_PATH")
	if dbPath == "" {
		dbPath = "shorty.db"
	}

	// Connect to database
	if err := database.Connect(dbPath); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes will be added here in future PRs
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "shorty",
			})
		})
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Shorty server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
