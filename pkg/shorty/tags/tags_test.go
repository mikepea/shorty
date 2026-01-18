package tags

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

func createTestLink(t *testing.T, db *gorm.DB, groupID, userID uint, slug string) models.Link {
	link := models.Link{
		GroupID:     groupID,
		CreatedByID: userID,
		Slug:        slug,
		URL:         "https://example.com",
		Title:       "Test Link",
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

func TestListTags(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	link := createTestLink(t, db, group.ID, user.ID, "test-link")

	// Create tags and associate with link
	tag1 := models.Tag{Name: "golang"}
	tag2 := models.Tag{Name: "programming"}
	db.Create(&tag1)
	db.Create(&tag2)
	db.Model(&link).Association("Tags").Append(&tag1, &tag2)

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tags []TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tags)

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestListTagsByGroup(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group1 := createTestGroup(t, db, "Group 1", user.ID)
	group2 := createTestGroup(t, db, "Group 2", user.ID)

	link1 := createTestLink(t, db, group1.ID, user.ID, "link1")
	link2 := createTestLink(t, db, group2.ID, user.ID, "link2")

	tag1 := models.Tag{Name: "group1-tag"}
	tag2 := models.Tag{Name: "group2-tag"}
	db.Create(&tag1)
	db.Create(&tag2)
	db.Model(&link1).Association("Tags").Append(&tag1)
	db.Model(&link2).Association("Tags").Append(&tag2)

	req, _ := http.NewRequest("GET", "/api/groups/1/tags", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tags []TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tags)

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}
	if len(tags) > 0 && tags[0].Name != "group1-tag" {
		t.Errorf("Expected 'group1-tag', got %s", tags[0].Name)
	}
}

func TestGetLinkTags(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	link := createTestLink(t, db, group.ID, user.ID, "test-link")

	tag := models.Tag{Name: "test-tag"}
	db.Create(&tag)
	db.Model(&link).Association("Tags").Append(&tag)

	req, _ := http.NewRequest("GET", "/api/links/test-link/tags", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tags []TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tags)

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}
}

func TestSetLinkTags(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	createTestLink(t, db, group.ID, user.ID, "test-link")

	body := SetTagsRequest{
		Tags: []string{"tag1", "tag2", "tag3"},
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/api/links/test-link/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tags []TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tags)

	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
}

func TestSetLinkTagsReplacesExisting(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	link := createTestLink(t, db, group.ID, user.ID, "test-link")

	// Add initial tag
	oldTag := models.Tag{Name: "old-tag"}
	db.Create(&oldTag)
	db.Model(&link).Association("Tags").Append(&oldTag)

	// Replace with new tags
	body := SetTagsRequest{
		Tags: []string{"new-tag"},
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/api/links/test-link/tags", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tags []TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tags)

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}
	if len(tags) > 0 && tags[0].Name != "new-tag" {
		t.Errorf("Expected 'new-tag', got %s", tags[0].Name)
	}
}

func TestAddLinkTag(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	createTestLink(t, db, group.ID, user.ID, "test-link")

	req, _ := http.NewRequest("POST", "/api/links/test-link/tags/new-tag", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var tag TagResponse
	json.Unmarshal(resp.Body.Bytes(), &tag)

	if tag.Name != "new-tag" {
		t.Errorf("Expected 'new-tag', got %s", tag.Name)
	}
}

func TestRemoveLinkTag(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	group := createTestGroup(t, db, "Test Group", user.ID)
	link := createTestLink(t, db, group.ID, user.ID, "test-link")

	tag := models.Tag{Name: "remove-me"}
	db.Create(&tag)
	db.Model(&link).Association("Tags").Append(&tag)

	req, _ := http.NewRequest("DELETE", "/api/links/test-link/tags/remove-me", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify tag was removed
	var updatedLink models.Link
	db.Preload("Tags").First(&updatedLink, link.ID)
	if len(updatedLink.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(updatedLink.Tags))
	}
}

func TestTagAccessControl(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	owner := createTestUser(t, db, "owner@example.com")
	other := createTestUser(t, db, "other@example.com")
	group := createTestGroup(t, db, "Test Group", owner.ID)
	createTestLink(t, db, group.ID, owner.ID, "test-link")

	// Other user should not be able to modify tags
	req, _ := http.NewRequest("POST", "/api/links/test-link/tags/hack", nil)
	req.Header.Set("Authorization", getAuthHeader(other))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}
