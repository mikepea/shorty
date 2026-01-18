package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/database"
	"github.com/mikepea/shorty/pkg/shorty/groups"
	"github.com/mikepea/shorty/pkg/shorty/links"
	"github.com/mikepea/shorty/pkg/shorty/models"
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

	// Run auto-migrations
	if err := models.AutoMigrate(database.GetDB()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Set up Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "shorty",
			})
		})

		// Auth routes
		authHandler := auth.NewHandler(database.GetDB())
		authHandler.RegisterRoutes(api.Group("/auth"))

		// Groups routes (protected)
		groupsHandler := groups.NewHandler(database.GetDB())
		groupsGroup := api.Group("/groups")
		groupsGroup.Use(auth.AuthMiddleware())
		groupsHandler.RegisterRoutes(groupsGroup)
		groupsHandler.RegisterMemberRoutes(groupsGroup)

		// Links routes (protected)
		linksHandler := links.NewHandler(database.GetDB())
		linksHandler.RegisterRoutes(api.Group("", auth.AuthMiddleware()))
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
