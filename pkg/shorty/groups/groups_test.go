package groups

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

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(db)

	groups := r.Group("/groups")
	groups.Use(auth.AuthMiddleware())
	handler.RegisterRoutes(groups)
	handler.RegisterMemberRoutes(groups)

	return r
}

func getAuthHeader(user models.User) string {
	token, _ := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	return "Bearer " + token
}

func TestCreateGroup(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	body := CreateGroupRequest{
		Name:        "Test Group",
		Description: "A test group",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response GroupResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Test Group" {
		t.Errorf("Expected name 'Test Group', got %s", response.Name)
	}
	if response.Role != "admin" {
		t.Errorf("Expected role 'admin', got %s", response.Role)
	}
}

func TestListGroups(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create a group
	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleMember,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/groups", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var groups []GroupResponse
	json.Unmarshal(resp.Body.Bytes(), &groups)

	if len(groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(groups))
	}
}

func TestGetGroup(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	group := models.Group{Name: "Test Group", Description: "Test description"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/groups/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response GroupResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Test Group" {
		t.Errorf("Expected name 'Test Group', got %s", response.Name)
	}
}

func TestGetGroupNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create group without adding user as member
	group := models.Group{Name: "Test Group"}
	db.Create(&group)

	req, _ := http.NewRequest("GET", "/groups/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestUpdateGroup(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	}
	db.Create(&membership)

	body := UpdateGroupRequest{Name: "Updated Group"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/groups/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response GroupResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Updated Group" {
		t.Errorf("Expected name 'Updated Group', got %s", response.Name)
	}
}

func TestUpdateGroupNotAdmin(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleMember, // Not admin
	}
	db.Create(&membership)

	body := UpdateGroupRequest{Name: "Updated Group"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/groups/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.Code)
	}
}

func TestListMembers(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/groups/1/members", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var members []MemberResponse
	json.Unmarshal(resp.Body.Bytes(), &members)

	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}
}

func TestAddMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	admin := createTestUser(t, db, "admin@example.com")
	newUser := createTestUser(t, db, "new@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	membership := models.GroupMembership{
		UserID:  admin.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	}
	db.Create(&membership)

	body := AddMemberRequest{
		Email: newUser.Email,
		Role:  "member",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response MemberResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Email != newUser.Email {
		t.Errorf("Expected email %s, got %s", newUser.Email, response.Email)
	}
}

func TestRemoveMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	admin := createTestUser(t, db, "admin@example.com")
	member := createTestUser(t, db, "member@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	db.Create(&models.GroupMembership{
		UserID:  admin.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	})
	db.Create(&models.GroupMembership{
		UserID:  member.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleMember,
	})

	req, _ := http.NewRequest("DELETE", "/groups/1/members/2", nil)
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCannotRemoveLastAdmin(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	admin := createTestUser(t, db, "admin@example.com")

	group := models.Group{Name: "Test Group"}
	db.Create(&group)
	db.Create(&models.GroupMembership{
		UserID:  admin.ID,
		GroupID: group.ID,
		Role:    models.GroupRoleAdmin,
	})

	// Try to remove self (last admin)
	req, _ := http.NewRequest("DELETE", "/groups/1/members/1", nil)
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", resp.Code, resp.Body.String())
	}
}
