package links

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
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

func createTestUser(t *testing.T, db *gorm.DB, email string) models.User {
	hash, _ := auth.HashPassword("password123")
	user := models.User{
		Email:        email,
		PasswordHash: hash,
		Name:         "Test User",
		SystemRole:   models.SystemRoleUser,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestGroup(t *testing.T, db *gorm.DB, name string, userID uint) models.Group {
	group := models.Group{Name: name}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("Failed to create test group: %v", err)
	}
	membership := models.GroupMembership{
		UserID:  userID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	}
	if err := db.Create(&membership).Error; err != nil {
		t.Fatalf("Failed to create test membership: %v", err)
	}
	return group
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(db)

	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	handler.RegisterRoutes(api)

	return r
}

func getAuthHeader(user models.User) string {
	token, _ := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	return "Bearer " + token
}

func TestCreateLink(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	body := CreateLinkRequest{
		URL:         "https://example.com",
		Slug:        "my-link",
		Title:       "Example",
		Description: "An example link",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/groups/1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Slug != "my-link" {
		t.Errorf("Expected slug 'my-link', got %s", response.Slug)
	}
	if response.GroupID != group.ID {
		t.Errorf("Expected group ID %d, got %d", group.ID, response.GroupID)
	}
}

func TestCreateLinkAutoSlug(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	createTestGroup(t, db, "Test Group", user.ID)

	body := CreateLinkRequest{
		URL:   "https://example.com",
		Title: "Example",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/groups/1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Slug == "" {
		t.Error("Expected auto-generated slug")
	}
	if len(response.Slug) != 8 {
		t.Errorf("Expected 8-char slug, got %d chars", len(response.Slug))
	}
}

func TestCreateLinkDuplicateSlug(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	createTestGroup(t, db, "Test Group", user.ID)

	// Create first link
	link := models.Link{
		GroupID:     1,
		CreatedByID: user.ID,
		Slug:        "existing-slug",
		URL:         "https://example.com",
	}
	db.Create(&link)

	// Try to create second link with same slug
	body := CreateLinkRequest{
		URL:  "https://another.com",
		Slug: "existing-slug",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/groups/1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestCreateLinkReservedSlug(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	createTestGroup(t, db, "Test Group", user.ID)

	body := CreateLinkRequest{
		URL:  "https://example.com",
		Slug: "api",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/groups/1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.Code)
	}
}

func TestListLinks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	// Create some links
	for i := 0; i < 3; i++ {
		db.Create(&models.Link{
			GroupID:     group.ID,
			CreatedByID: user.ID,
			Slug:        generateRandomString(8, "abcdef"),
			URL:         "https://example.com",
		})
	}

	req, _ := http.NewRequest("GET", "/api/groups/1/links", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var links []LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &links)

	if len(links) != 3 {
		t.Errorf("Expected 3 links, got %d", len(links))
	}
}

func TestGetLinkBySlug(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "test-link",
		URL:         "https://example.com",
		Title:       "Test Link",
	}
	db.Create(&link)

	req, _ := http.NewRequest("GET", "/api/links/test-link", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Title != "Test Link" {
		t.Errorf("Expected title 'Test Link', got %s", response.Title)
	}
}

func TestGetPrivateLinkNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	owner := createTestUser(t, db, "owner@example.com")
	other := createTestUser(t, db, "other@example.com")
	group := createTestGroup(t, db, "Test Group", owner.ID)

	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: owner.ID,
		Slug:        "private-link",
		URL:         "https://example.com",
		IsPublic:    false,
	}
	db.Create(&link)

	req, _ := http.NewRequest("GET", "/api/links/private-link", nil)
	req.Header.Set("Authorization", getAuthHeader(other))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestGetPublicLinkNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	owner := createTestUser(t, db, "owner@example.com")
	other := createTestUser(t, db, "other@example.com")
	group := createTestGroup(t, db, "Test Group", owner.ID)

	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: owner.ID,
		Slug:        "public-link",
		URL:         "https://example.com",
		IsPublic:    true,
	}
	db.Create(&link)

	req, _ := http.NewRequest("GET", "/api/links/public-link", nil)
	req.Header.Set("Authorization", getAuthHeader(other))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

func TestUpdateLink(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "test-link",
		URL:         "https://example.com",
		Title:       "Old Title",
	}
	db.Create(&link)

	body := UpdateLinkRequest{
		Title: "New Title",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/api/links/test-link", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Title != "New Title" {
		t.Errorf("Expected title 'New Title', got %s", response.Title)
	}
}

func TestDeleteLink(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "test-link",
		URL:         "https://example.com",
	}
	db.Create(&link)

	req, _ := http.NewRequest("DELETE", "/api/links/test-link", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestSearchLinks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	// Create links with different titles
	db.Create(&models.Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "link1",
		URL:         "https://example.com",
		Title:       "Golang Tutorial",
	})
	db.Create(&models.Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "link2",
		URL:         "https://example.com",
		Title:       "Python Guide",
	})

	req, _ := http.NewRequest("GET", "/api/links?q=Golang", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var links []LinkResponse
	json.Unmarshal(resp.Body.Bytes(), &links)

	if len(links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(links))
	}
	if links[0].Title != "Golang Tutorial" {
		t.Errorf("Expected 'Golang Tutorial', got %s", links[0].Title)
	}
}
