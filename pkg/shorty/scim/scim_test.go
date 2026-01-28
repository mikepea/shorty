package scim

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

// createGlobalOrg creates the global organization for tests
func createGlobalOrg(t *testing.T, db *gorm.DB) *models.Organization {
	org := &models.Organization{
		Name:     "Shorty Global",
		Slug:     "shorty-global",
		IsGlobal: true,
	}
	if err := db.Create(org).Error; err != nil {
		t.Fatalf("Failed to create global organization: %v", err)
	}
	return org
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func createTestUser(t *testing.T, db *gorm.DB, email, name string) *models.User {
	user := &models.User{
		Email:      email,
		Name:       name,
		Active:     true,
		SystemRole: models.SystemRoleUser,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestGroup(t *testing.T, db *gorm.DB, name string) *models.Group {
	group := &models.Group{
		Name: name,
	}
	if err := db.Create(group).Error; err != nil {
		t.Fatalf("Failed to create test group: %v", err)
	}
	return group
}

// User Tests

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	createTestUser(t, db, "user1@test.com", "User One")
	createTestUser(t, db, "user2@test.com", "User Two")

	r.GET("/scim/v2/Users", h.ListUsers)

	req := httptest.NewRequest("GET", "/scim/v2/Users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.TotalResults != 2 {
		t.Errorf("Expected 2 users, got %d", resp.TotalResults)
	}
}

func TestListUsersWithFilter(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	createTestUser(t, db, "john@test.com", "John Doe")
	createTestUser(t, db, "jane@test.com", "Jane Doe")

	r.GET("/scim/v2/Users", h.ListUsers)

	req := httptest.NewRequest("GET", "/scim/v2/Users?filter=userName%20eq%20%22john@test.com%22", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.TotalResults != 1 {
		t.Errorf("Expected 1 user matching filter, got %d", resp.TotalResults)
	}
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	r.POST("/scim/v2/Users", h.CreateUser)

	body := CreateUserRequest{
		Schemas:  []string{SchemaUser},
		UserName: "newuser@test.com",
		Name: Name{
			GivenName:  "New",
			FamilyName: "User",
		},
		DisplayName: "New User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/scim/v2/Users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var user User
	json.Unmarshal(w.Body.Bytes(), &user)

	if user.UserName != "newuser@test.com" {
		t.Errorf("Expected userName newuser@test.com, got %s", user.UserName)
	}
}

func TestGetUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	user := createTestUser(t, db, "test@test.com", "Test User")

	r.GET("/scim/v2/Users/:id", h.GetUser)

	req := httptest.NewRequest("GET", "/scim/v2/Users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp User
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.UserName != user.Email {
		t.Errorf("Expected userName %s, got %s", user.Email, resp.UserName)
	}
}

func TestPatchUserActive(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	createTestUser(t, db, "test@test.com", "Test User")

	r.PATCH("/scim/v2/Users/:id", h.PatchUser)

	patch := PatchOp{
		Schemas: []string{SchemaPatchOp},
		Operations: []PatchOperation{
			{
				Op:    "replace",
				Path:  "active",
				Value: false,
			},
		},
	}
	jsonBody, _ := json.Marshal(patch)

	req := httptest.NewRequest("PATCH", "/scim/v2/Users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var user User
	json.Unmarshal(w.Body.Bytes(), &user)

	if user.Active != false {
		t.Errorf("Expected active to be false")
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewUserHandler(db, "http://localhost:8080")

	createTestUser(t, db, "test@test.com", "Test User")

	r.DELETE("/scim/v2/Users/:id", h.DeleteUser)

	req := httptest.NewRequest("DELETE", "/scim/v2/Users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Verify user is deleted
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected user to be deleted")
	}
}

// Group Tests

func TestListGroups(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewGroupHandler(db, "http://localhost:8080")

	createTestGroup(t, db, "Group One")
	createTestGroup(t, db, "Group Two")

	r.GET("/scim/v2/Groups", h.ListGroups)

	req := httptest.NewRequest("GET", "/scim/v2/Groups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.TotalResults != 2 {
		t.Errorf("Expected 2 groups, got %d", resp.TotalResults)
	}
}

func TestCreateGroup(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewGroupHandler(db, "http://localhost:8080")

	r.POST("/scim/v2/Groups", h.CreateGroup)

	body := CreateGroupRequest{
		Schemas:     []string{SchemaGroup},
		DisplayName: "New Group",
		ExternalID:  "ext-123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/scim/v2/Groups", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var group Group
	json.Unmarshal(w.Body.Bytes(), &group)

	if group.DisplayName != "New Group" {
		t.Errorf("Expected displayName 'New Group', got %s", group.DisplayName)
	}
}

func TestPatchGroupMembers(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	gh := NewGroupHandler(db, "http://localhost:8080")

	user := createTestUser(t, db, "test@test.com", "Test User")
	group := createTestGroup(t, db, "Test Group")

	r.PATCH("/scim/v2/Groups/:id", gh.PatchGroup)

	patch := PatchOp{
		Schemas: []string{SchemaPatchOp},
		Operations: []PatchOperation{
			{
				Op:   "add",
				Path: "members",
				Value: []map[string]interface{}{
					{"value": "1"},
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(patch)

	req := httptest.NewRequest("PATCH", "/scim/v2/Groups/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify membership
	var membership models.GroupMembership
	err := db.Where("user_id = ? AND group_id = ?", user.ID, group.ID).First(&membership).Error
	if err != nil {
		t.Errorf("Expected membership to be created: %v", err)
	}
}

// Service Provider Config Test

func TestGetServiceProviderConfig(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	h := NewConfigHandler(db, "http://localhost:8080")

	r.GET("/scim/v2/ServiceProviderConfig", h.GetServiceProviderConfig)

	req := httptest.NewRequest("GET", "/scim/v2/ServiceProviderConfig", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var config ServiceProviderConfig
	json.Unmarshal(w.Body.Bytes(), &config)

	if !config.Patch.Supported {
		t.Error("Expected patch to be supported")
	}
	if !config.Filter.Supported {
		t.Error("Expected filter to be supported")
	}
}

// SCIM Token Test

func TestSCIMTokenGeneration(t *testing.T) {
	db := setupTestDB(t)
	org := createGlobalOrg(t, db)

	token, scimToken, err := GenerateSCIMToken(db, org.ID, "Test Token")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if len(token) != 64 { // 32 bytes hex encoded
		t.Errorf("Expected token length 64, got %d", len(token))
	}

	if scimToken.TokenPrefix != token[:8] {
		t.Errorf("Expected token prefix %s, got %s", token[:8], scimToken.TokenPrefix)
	}

	if scimToken.OrganizationID != org.ID {
		t.Errorf("Expected organization ID %d, got %d", org.ID, scimToken.OrganizationID)
	}

	// Validate the token
	validatedToken, err := ValidateSCIMToken(db, token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if validatedToken.ID != scimToken.ID {
		t.Errorf("Expected token ID %d, got %d", scimToken.ID, validatedToken.ID)
	}
}

func TestSCIMAuthMiddleware(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter()
	org := createGlobalOrg(t, db)

	token, _, _ := GenerateSCIMToken(db, org.ID, "Test Token")

	r.Use(SCIMAuthMiddleware(db))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Test without auth
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 without auth, got %d", w.Code)
	}

	// Test with valid token
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 with valid token, got %d", w.Code)
	}

	// Test with invalid token
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 with invalid token, got %d", w.Code)
	}
}
