package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/apikeys"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/groups"
	"github.com/mikepea/shorty/pkg/shorty/importexport"
	"github.com/mikepea/shorty/pkg/shorty/links"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"github.com/mikepea/shorty/pkg/shorty/redirect"
	"github.com/mikepea/shorty/pkg/shorty/tags"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	if err := models.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create the global organization (required for redirect functionality)
	globalOrg := models.Organization{
		Name:     "Shorty Global",
		Slug:     "shorty-global",
		IsGlobal: true,
	}
	if err := db.Create(&globalOrg).Error; err != nil {
		t.Fatalf("Failed to create global organization: %v", err)
	}

	return db
}

// setupFullServer creates a Gin engine with all routes registered
// This mirrors the setup in cmd/shorty-server/main.go
func setupFullServer(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
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
		authHandler := auth.NewHandler(db)
		authHandler.RegisterRoutes(api.Group("/auth"))

		// Combined auth middleware (accepts JWT or API key)
		combinedAuth := apikeys.CombinedAuthMiddleware(db)

		// API keys routes (JWT only - need to be logged in to manage keys)
		apiKeysHandler := apikeys.NewHandler(db)
		apiKeysHandler.RegisterRoutes(api.Group("", auth.AuthMiddleware()))

		// Groups routes (protected - accepts JWT or API key)
		groupsHandler := groups.NewHandler(db)
		groupsGroup := api.Group("/groups")
		groupsGroup.Use(combinedAuth)
		groupsHandler.RegisterRoutes(groupsGroup)
		groupsHandler.RegisterMemberRoutes(groupsGroup)

		// Links routes (protected - accepts JWT or API key)
		linksHandler := links.NewHandler(db)
		linksHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Tags routes (protected - accepts JWT or API key)
		tagsHandler := tags.NewHandler(db)
		tagsHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Import/Export routes (protected - accepts JWT or API key)
		importExportHandler := importexport.NewHandler(db)
		importExportHandler.RegisterRoutes(api.Group("", combinedAuth))
	}

	// Redirect routes (public, must be registered LAST to avoid conflicts)
	redirectHandler := redirect.NewHandler(db)
	redirectHandler.RegisterRoutes(r)

	return r
}

// TestServerStartup verifies that all routes can be registered without conflicts
// This test would fail if there are route parameter conflicts (like :id vs :groupId)
func TestServerStartup(t *testing.T) {
	db := setupTestDB(t)

	// This will panic if there are route conflicts
	router := setupFullServer(db)

	if router == nil {
		t.Fatal("Expected router to be created")
	}
}

// TestHealthEndpoint verifies the health endpoint responds correctly
func TestHealthEndpoint(t *testing.T) {
	db := setupTestDB(t)
	router := setupFullServer(db)

	req, _ := http.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

// TestAPIHealthEndpoint verifies the API health endpoint responds correctly
func TestAPIHealthEndpoint(t *testing.T) {
	db := setupTestDB(t)
	router := setupFullServer(db)

	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

// TestProtectedEndpointsRequireAuth verifies that protected endpoints return 401 without auth
func TestProtectedEndpointsRequireAuth(t *testing.T) {
	db := setupTestDB(t)
	router := setupFullServer(db)

	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/groups"},
		{"POST", "/api/groups"},
		{"GET", "/api/links"},
		{"GET", "/api/tags"},
	}

	for _, endpoint := range protectedEndpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req, _ := http.NewRequest(endpoint.method, endpoint.path, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401 for %s %s, got %d", endpoint.method, endpoint.path, resp.Code)
			}
		})
	}
}

// TestPublicEndpointsNoAuth verifies that public endpoints don't require auth
func TestPublicEndpointsNoAuth(t *testing.T) {
	db := setupTestDB(t)
	router := setupFullServer(db)

	publicEndpoints := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/api/health", http.StatusOK},
		{"POST", "/api/auth/register", http.StatusBadRequest}, // Bad request (no body), but not 401
		{"POST", "/api/auth/login", http.StatusBadRequest},    // Bad request (no body), but not 401
		{"GET", "/nonexistent-slug", http.StatusNotFound},     // 404 for missing link, but not 401
	}

	for _, endpoint := range publicEndpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req, _ := http.NewRequest(endpoint.method, endpoint.path, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != endpoint.expectedCode {
				t.Errorf("Expected status %d for %s %s, got %d", endpoint.expectedCode, endpoint.method, endpoint.path, resp.Code)
			}
		})
	}
}
