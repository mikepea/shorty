package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/admin"
	"github.com/mikepea/shorty/pkg/shorty/apikeys"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/database"
	"github.com/mikepea/shorty/pkg/shorty/groups"
	"github.com/mikepea/shorty/pkg/shorty/importexport"
	"github.com/mikepea/shorty/pkg/shorty/links"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"github.com/mikepea/shorty/pkg/shorty/redirect"
	"github.com/mikepea/shorty/pkg/shorty/tags"
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

	// Create default admin user if no admin exists
	if err := ensureAdminExists(); err != nil {
		log.Fatalf("Failed to ensure admin user exists: %v", err)
	}

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

		// Auth routes (public)
		authHandler := auth.NewHandler(database.GetDB())
		authHandler.RegisterRoutes(api.Group("/auth"))

		// Combined auth middleware (accepts JWT or API key)
		combinedAuth := apikeys.CombinedAuthMiddleware(database.GetDB())

		// API keys routes (JWT only - need to be logged in to manage keys)
		apiKeysHandler := apikeys.NewHandler(database.GetDB())
		apiKeysHandler.RegisterRoutes(api.Group("", auth.AuthMiddleware()))

		// Groups routes (protected - accepts JWT or API key)
		groupsHandler := groups.NewHandler(database.GetDB())
		groupsGroup := api.Group("/groups")
		groupsGroup.Use(combinedAuth)
		groupsHandler.RegisterRoutes(groupsGroup)
		groupsHandler.RegisterMemberRoutes(groupsGroup)

		// Links routes (protected - accepts JWT or API key)
		linksHandler := links.NewHandler(database.GetDB())
		linksHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Tags routes (protected - accepts JWT or API key)
		tagsHandler := tags.NewHandler(database.GetDB())
		tagsHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Import/Export routes (protected - accepts JWT or API key)
		importExportHandler := importexport.NewHandler(database.GetDB())
		importExportHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Admin routes (JWT only, admin role required)
		adminHandler := admin.NewHandler(database.GetDB())
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(), auth.RequireAdmin())
		adminHandler.RegisterRoutes(adminGroup)
	}

	// Redirect routes (public, must be registered LAST to avoid conflicts)
	redirectHandler := redirect.NewHandler(database.GetDB())
	redirectHandler.RegisterRoutes(r)

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

// ensureAdminExists creates a default admin user if no admin exists in the database
func ensureAdminExists() error {
	db := database.GetDB()

	// Check if any admin user exists
	var count int64
	if err := db.Model(&models.User{}).Where("system_role = ?", models.SystemRoleAdmin).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Admin already exists
	}

	// Create default admin user
	hashedPassword, err := auth.HashPassword("changeme")
	if err != nil {
		return err
	}

	adminUser := models.User{
		Email:        "admin@localhost",
		Name:         "Admin",
		PasswordHash: hashedPassword,
		SystemRole:   models.SystemRoleAdmin,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	// Create personal group for admin
	personalGroup := models.Group{
		Name:        "Admin's Links",
		Description: "Personal links for Admin",
	}
	if err := db.Create(&personalGroup).Error; err != nil {
		return err
	}

	// Add admin as admin of personal group
	membership := models.GroupMembership{
		UserID:  adminUser.ID,
		GroupID: personalGroup.ID,
		Role:    models.GroupRoleAdmin,
	}
	if err := db.Create(&membership).Error; err != nil {
		return err
	}

	log.Printf("Created default admin user: admin@localhost (password: changeme)")
	return nil
}
