package importexport

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

func createTestLink(t *testing.T, db *gorm.DB, groupID, userID uint, slug, url string) models.Link {
	link := models.Link{
		GroupID:     groupID,
		CreatedByID: userID,
		Slug:        slug,
		URL:         url,
		Title:       "Test Link",
		Description: "Test description",
		IsPublic:    true,
		IsUnread:    false,
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

	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	handler.RegisterRoutes(api)

	return r
}

func getAuthHeader(user models.User) string {
	token, _ := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	return "Bearer " + token
}

func TestImportBookmarks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	req := ImportRequest{
		GroupID: group.ID,
		Bookmarks: []PinboardBookmark{
			{
				Href:        "https://example.com",
				Description: "Example Site",
				Extended:    "This is an example",
				Tags:        "test example",
				Time:        "2024-01-15T10:30:00Z",
				Shared:      "yes",
				ToRead:      "no",
			},
			{
				Href:        "https://golang.org",
				Description: "Go Programming",
				Extended:    "The Go programming language",
				Tags:        "golang programming",
				Time:        "2024-01-16T14:00:00Z",
				Shared:      "no",
				ToRead:      "yes",
			},
		},
	}
	jsonBody, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", "/api/import", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result ImportResult
	json.Unmarshal(resp.Body.Bytes(), &result)

	if result.Imported != 2 {
		t.Errorf("Expected 2 imported, got %d", result.Imported)
	}

	if result.Skipped != 0 {
		t.Errorf("Expected 0 skipped, got %d", result.Skipped)
	}

	// Verify links were created
	var count int64
	db.Model(&models.Link{}).Where("group_id = ?", group.ID).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 links in database, got %d", count)
	}

	// Verify tags were created
	var tagCount int64
	db.Model(&models.Tag{}).Count(&tagCount)
	if tagCount != 4 {
		t.Errorf("Expected 4 tags, got %d", tagCount)
	}
}

func TestImportBookmarksNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	otherUser := createTestUser(t, db, "other@example.com")
	group := createTestGroup(t, db, "Test Group", otherUser.ID)

	req := ImportRequest{
		GroupID: group.ID,
		Bookmarks: []PinboardBookmark{
			{
				Href:        "https://example.com",
				Description: "Example Site",
			},
		},
	}
	jsonBody, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", "/api/import", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestExportBookmarks(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	// Create links with tags
	link1 := createTestLink(t, db, group.ID, user.ID, "link1", "https://example.com")
	link2 := createTestLink(t, db, group.ID, user.ID, "link2", "https://golang.org")

	tag1 := models.Tag{Name: "test"}
	tag2 := models.Tag{Name: "golang"}
	db.Create(&tag1)
	db.Create(&tag2)
	db.Model(&link1).Association("Tags").Append(&tag1)
	db.Model(&link2).Association("Tags").Append(&tag2)

	httpReq, _ := http.NewRequest("GET", "/api/export", nil)
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var bookmarks []ExportBookmark
	json.Unmarshal(resp.Body.Bytes(), &bookmarks)

	if len(bookmarks) != 2 {
		t.Errorf("Expected 2 bookmarks, got %d", len(bookmarks))
	}
}

func TestExportBookmarksByGroup(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group1 := createTestGroup(t, db, "Group 1", user.ID)
	group2 := createTestGroup(t, db, "Group 2", user.ID)

	createTestLink(t, db, group1.ID, user.ID, "link1", "https://example.com")
	createTestLink(t, db, group2.ID, user.ID, "link2", "https://golang.org")

	// Export only from group1
	httpReq, _ := http.NewRequest("GET", "/api/export?group_id=1", nil)
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var bookmarks []ExportBookmark
	json.Unmarshal(resp.Body.Bytes(), &bookmarks)

	if len(bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(bookmarks))
	}

	if bookmarks[0].Href != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", bookmarks[0].Href)
	}
}

func TestExportSingleBookmark(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	link := createTestLink(t, db, group.ID, user.ID, "test-link", "https://example.com")
	tag := models.Tag{Name: "test"}
	db.Create(&tag)
	db.Model(&link).Association("Tags").Append(&tag)

	httpReq, _ := http.NewRequest("GET", "/api/export/test-link", nil)
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var bookmark ExportBookmark
	json.Unmarshal(resp.Body.Bytes(), &bookmark)

	if bookmark.Href != "https://example.com" {
		t.Errorf("Expected https://example.com, got %s", bookmark.Href)
	}

	if bookmark.Tags != "test" {
		t.Errorf("Expected 'test' tag, got %s", bookmark.Tags)
	}
}

func TestExportSingleBookmarkNotFound(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	httpReq, _ := http.NewRequest("GET", "/api/export/nonexistent", nil)
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestExportPrivateLinkNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	otherUser := createTestUser(t, db, "other@example.com")
	group := createTestGroup(t, db, "Test Group", otherUser.ID)

	// Create private link
	link := models.Link{
		GroupID:     group.ID,
		CreatedByID: otherUser.ID,
		Slug:        "private-link",
		URL:         "https://secret.example.com",
		Title:       "Private Link",
		IsPublic:    false,
	}
	db.Create(&link)

	httpReq, _ := http.NewRequest("GET", "/api/export/private-link", nil)
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestImportPreservesTimestamp(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	req := ImportRequest{
		GroupID: group.ID,
		Bookmarks: []PinboardBookmark{
			{
				Href:        "https://example.com",
				Description: "Example Site",
				Time:        "2020-06-15T10:30:00Z",
				Shared:      "yes",
				ToRead:      "no",
			},
		},
	}
	jsonBody, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", "/api/import", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify timestamp was preserved
	var link models.Link
	db.Where("group_id = ?", group.ID).First(&link)

	expectedYear := 2020
	if link.CreatedAt.Year() != expectedYear {
		t.Errorf("Expected year %d, got %d", expectedYear, link.CreatedAt.Year())
	}
}

func TestImportVisibilityAndUnread(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)

	req := ImportRequest{
		GroupID: group.ID,
		Bookmarks: []PinboardBookmark{
			{
				Href:        "https://public.example.com",
				Description: "Public Read",
				Shared:      "yes",
				ToRead:      "no",
			},
			{
				Href:        "https://private.example.com",
				Description: "Private Unread",
				Shared:      "no",
				ToRead:      "yes",
			},
		},
	}
	jsonBody, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", "/api/import", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, httpReq)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify first link is public and read
	var publicLink models.Link
	db.Where("url = ?", "https://public.example.com").First(&publicLink)
	if !publicLink.IsPublic {
		t.Error("Expected public link to be public")
	}
	if publicLink.IsUnread {
		t.Error("Expected public link to not be unread")
	}

	// Verify second link is private and unread
	var privateLink models.Link
	db.Where("url = ?", "https://private.example.com").First(&privateLink)
	if privateLink.IsPublic {
		t.Error("Expected private link to not be public")
	}
	if !privateLink.IsUnread {
		t.Error("Expected private link to be unread")
	}
}
