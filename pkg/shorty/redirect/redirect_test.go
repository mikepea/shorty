package redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	models.AutoMigrate(db)
	return db
}

func createTestLink(t *testing.T, db *gorm.DB, slug, url string, isPublic bool) models.Link {
	link := models.Link{
		GroupID:     1,
		CreatedByID: 1,
		Slug:        slug,
		URL:         url,
		Title:       "Test Link",
		IsPublic:    isPublic,
		ClickCount:  0,
	}
	if err := db.Create(&link).Error; err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}
	return link
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(db)
	handler.RegisterRoutes(r)
	return r
}

func TestRedirectPublicLink(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	createTestLink(t, db, "test-link", "https://example.com", true)

	req, _ := http.NewRequest("GET", "/test-link", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.Code)
	}

	location := resp.Header().Get("Location")
	if location != "https://example.com" {
		t.Errorf("Expected Location 'https://example.com', got %s", location)
	}
}

func TestRedirectPrivateLink(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	createTestLink(t, db, "private-link", "https://secret.example.com", false)

	req, _ := http.NewRequest("GET", "/private-link", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Private links still redirect - the URL itself isn't secret
	if resp.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.Code)
	}

	location := resp.Header().Get("Location")
	if location != "https://secret.example.com" {
		t.Errorf("Expected Location 'https://secret.example.com', got %s", location)
	}
}

func TestRedirectNotFound(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestRedirectIncrementsClickCount(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	link := createTestLink(t, db, "click-test", "https://example.com", true)

	// Initial click count should be 0
	if link.ClickCount != 0 {
		t.Errorf("Expected initial click count 0, got %d", link.ClickCount)
	}

	// Make a request
	req, _ := http.NewRequest("GET", "/click-test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.Code)
	}

	// Wait a bit for the goroutine to complete
	time.Sleep(100 * time.Millisecond)

	// Check click count was incremented
	var updatedLink models.Link
	db.First(&updatedLink, link.ID)
	if updatedLink.ClickCount != 1 {
		t.Errorf("Expected click count 1, got %d", updatedLink.ClickCount)
	}

	// Make another request
	req, _ = http.NewRequest("GET", "/click-test", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Wait for goroutine
	time.Sleep(100 * time.Millisecond)

	// Check click count again
	db.First(&updatedLink, link.ID)
	if updatedLink.ClickCount != 2 {
		t.Errorf("Expected click count 2, got %d", updatedLink.ClickCount)
	}
}

func TestRedirectWithQueryParams(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	createTestLink(t, db, "query-test", "https://example.com/page", true)

	// Query params in the request should not affect the redirect
	// The redirect goes to the stored URL
	req, _ := http.NewRequest("GET", "/query-test?foo=bar", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", resp.Code)
	}

	location := resp.Header().Get("Location")
	if location != "https://example.com/page" {
		t.Errorf("Expected Location 'https://example.com/page', got %s", location)
	}
}
