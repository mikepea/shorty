package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/admin"
	"github.com/mikepea/shorty/pkg/shorty/apikeys"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/database"
	"github.com/mikepea/shorty/pkg/shorty/groups"
	"github.com/mikepea/shorty/pkg/shorty/importexport"
	"github.com/mikepea/shorty/pkg/shorty/links"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"github.com/mikepea/shorty/pkg/shorty/oidc"
	"github.com/mikepea/shorty/pkg/shorty/redirect"
	"github.com/mikepea/shorty/pkg/shorty/scim"
	"github.com/mikepea/shorty/pkg/shorty/tags"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/mikepea/shorty/api/swagger"
)

// @title Shorty API
// @version 1.0
// @description A modern URL shortener with team collaboration, SSO, and SCIM support.

// @contact.name Shorty Support
// @contact.url https://github.com/mikepea/shorty

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token or API key. Format: "Bearer {token}"

// @securityDefinitions.apikey SCIMAuth
// @in header
// @name Authorization
// @description SCIM bearer token. Format: "Bearer {scim_token}"

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

	// Ensure the global organization exists (must run before admin creation)
	globalOrg, err := ensureGlobalOrgExists()
	if err != nil {
		log.Fatalf("Failed to ensure global organization exists: %v", err)
	}

	// Create default admin user if no admin exists
	if err := ensureAdminExists(globalOrg); err != nil {
		log.Fatalf("Failed to ensure admin user exists: %v", err)
	}

	// Get base URL from environment or use default
	baseURL := os.Getenv("SHORTY_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Set up Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

		// OIDC routes
		oidcHandler := oidc.NewHandler(database.GetDB(), baseURL)
		oidcHandler.RegisterRoutes(api.Group("/oidc"))
		oidcHandler.RegisterAdminRoutes(adminGroup.Group("/oidc"))

		// SCIM token management (admin only)
		scimTokenHandler := scim.NewTokenHandler(database.GetDB())
		scimTokenHandler.RegisterAdminRoutes(adminGroup)
	}

	// SCIM routes (bearer token auth, outside /api to follow SCIM spec)
	scimGroup := r.Group("/scim/v2")
	scimGroup.Use(scim.SCIMAuthMiddleware(database.GetDB()))
	{
		scimUserHandler := scim.NewUserHandler(database.GetDB(), baseURL)
		scimUserHandler.RegisterRoutes(scimGroup)

		scimGroupHandler := scim.NewGroupHandler(database.GetDB(), baseURL)
		scimGroupHandler.RegisterRoutes(scimGroup)

		scimConfigHandler := scim.NewConfigHandler(database.GetDB(), baseURL)
		scimConfigHandler.RegisterRoutes(scimGroup)
	}

	// Serve static frontend files if web/dist exists
	webDistPath := "./web/dist"
	if _, err := os.Stat(webDistPath); err == nil {
		// Serve static assets (JS, CSS, images, etc.)
		r.Static("/assets", filepath.Join(webDistPath, "assets"))

		// Serve other static files at root (favicon, etc.)
		r.StaticFile("/favicon.ico", filepath.Join(webDistPath, "favicon.ico"))
		r.StaticFile("/robots.txt", filepath.Join(webDistPath, "robots.txt"))

		// SPA fallback - serve index.html for frontend routes
		indexHTML := filepath.Join(webDistPath, "index.html")
		spaRoutes := []string{"/", "/login", "/register", "/dashboard", "/links", "/groups", "/settings", "/admin"}
		for _, route := range spaRoutes {
			route := route // capture loop variable
			r.GET(route, func(c *gin.Context) {
				c.File(indexHTML)
			})
		}
		// Also handle sub-routes like /links/:slug
		r.GET("/links/*path", func(c *gin.Context) {
			c.File(indexHTML)
		})
		r.GET("/groups/*path", func(c *gin.Context) {
			c.File(indexHTML)
		})

		log.Println("Serving frontend from ./web/dist")
	} else {
		log.Println("No frontend build found at ./web/dist - API only mode")
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

// ensureGlobalOrgExists creates the "Shorty Global" organization if it doesn't exist.
// This organization serves as the default for public signups and unrecognized domains.
// Returns the global organization.
func ensureGlobalOrgExists() (*models.Organization, error) {
	db := database.GetDB()

	// Check if global org already exists
	var globalOrg models.Organization
	err := db.Where("is_global = ?", true).First(&globalOrg).Error
	if err == nil {
		return &globalOrg, nil // Already exists
	}

	// Create the global organization
	globalOrg = models.Organization{
		Name:     "Shorty Global",
		Slug:     "shorty-global",
		IsGlobal: true,
	}

	if err := db.Create(&globalOrg).Error; err != nil {
		return nil, err
	}

	log.Printf("Created global organization: %s (ID: %d)", globalOrg.Name, globalOrg.ID)

	// Migrate any existing data to the global organization
	if err := migrateExistingDataToGlobalOrg(db, globalOrg.ID); err != nil {
		log.Printf("Warning: Error migrating existing data to global org: %v", err)
	}

	return &globalOrg, nil
}

// migrateExistingDataToGlobalOrg assigns any existing groups, links, OIDC providers,
// and SCIM tokens without an organization to the global organization.
// This handles upgrades from pre-multi-tenancy versions.
func migrateExistingDataToGlobalOrg(db *gorm.DB, globalOrgID uint) error {
	// Migrate groups without an organization
	if err := db.Model(&models.Group{}).Where("organization_id = 0 OR organization_id IS NULL").
		Update("organization_id", globalOrgID).Error; err != nil {
		return err
	}

	// Migrate links without an organization
	if err := db.Model(&models.Link{}).Where("organization_id = 0 OR organization_id IS NULL").
		Update("organization_id", globalOrgID).Error; err != nil {
		return err
	}

	// Migrate OIDC providers without an organization
	if err := db.Model(&models.OIDCProvider{}).Where("organization_id = 0 OR organization_id IS NULL").
		Update("organization_id", globalOrgID).Error; err != nil {
		return err
	}

	// Migrate SCIM tokens without an organization
	if err := db.Model(&models.SCIMToken{}).Where("organization_id = 0 OR organization_id IS NULL").
		Update("organization_id", globalOrgID).Error; err != nil {
		return err
	}

	return nil
}

// ensureAdminExists creates a default admin user if no admin exists in the database.
// The admin is added to the global organization.
func ensureAdminExists(globalOrg *models.Organization) error {
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
		Email:        "admin@shorty.local",
		Name:         "Admin",
		PasswordHash: hashedPassword,
		SystemRole:   models.SystemRoleAdmin,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	// Add admin to global organization as admin
	orgMembership := models.OrganizationMembership{
		OrganizationID: globalOrg.ID,
		UserID:         adminUser.ID,
		Role:           models.OrgRoleAdmin,
	}
	if err := db.Create(&orgMembership).Error; err != nil {
		return err
	}

	// Create personal group for admin within the global organization
	personalGroup := models.Group{
		OrganizationID: globalOrg.ID,
		Name:           "Admin's Links",
		Description:    "Personal links for Admin",
	}
	if err := db.Create(&personalGroup).Error; err != nil {
		return err
	}

	// Add admin as admin of personal group
	groupMembership := models.GroupMembership{
		UserID:  adminUser.ID,
		GroupID: personalGroup.ID,
		Role:    models.GroupRoleAdmin,
	}
	if err := db.Create(&groupMembership).Error; err != nil {
		return err
	}

	log.Printf("Created default admin user: admin@shorty.local (password: changeme)")
	return nil
}
