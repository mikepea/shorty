package admin

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
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func createTestUser(t *testing.T, db *gorm.DB, email, name string, role models.SystemRole) *models.User {
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
		SystemRole:   role,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestLink(t *testing.T, db *gorm.DB, createdByID, groupID uint, slug string) *models.Link {
	link := &models.Link{
		URL:         "https://example.com/" + slug,
		Slug:        slug,
		CreatedByID: createdByID,
		GroupID:     groupID,
	}
	if err := db.Create(link).Error; err != nil {
		t.Fatalf("Failed to create test link: %v", err)
	}
	return link
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

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	// Create admin and regular users
	admin := createTestUser(t, db, "admin@test.com", "Admin User", models.SystemRoleAdmin)
	createTestUser(t, db, "user1@test.com", "User One", models.SystemRoleUser)
	createTestUser(t, db, "user2@test.com", "User Two", models.SystemRoleUser)

	r.GET("/admin/users", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.ListUsers(c)
	})

	req := httptest.NewRequest("GET", "/admin/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var users []UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

func TestListUsersWithSearch(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin User", models.SystemRoleAdmin)
	createTestUser(t, db, "john@test.com", "John Doe", models.SystemRoleUser)
	createTestUser(t, db, "jane@test.com", "Jane Doe", models.SystemRoleUser)

	r.GET("/admin/users", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.ListUsers(c)
	})

	req := httptest.NewRequest("GET", "/admin/users?q=john", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var users []UserResponse
	json.Unmarshal(w.Body.Bytes(), &users)

	if len(users) != 1 {
		t.Errorf("Expected 1 user matching search, got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)
	user := createTestUser(t, db, "user@test.com", "Test User", models.SystemRoleUser)

	// Create some links for the user
	group := createTestGroup(t, db, "Test Group")
	createTestLink(t, db, user.ID, group.ID, "link1")
	createTestLink(t, db, user.ID, group.ID, "link2")

	// Add group membership
	db.Create(&models.GroupMembership{UserID: user.ID, GroupID: group.ID, Role: models.GroupRoleAdmin})

	r.GET("/admin/users/:id", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.GetUser(c)
	})

	req := httptest.NewRequest("GET", "/admin/users/"+string(rune('0'+user.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp UserResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, resp.Email)
	}
	if resp.LinkCount != 2 {
		t.Errorf("Expected 2 links, got %d", resp.LinkCount)
	}
	if resp.GroupCount != 1 {
		t.Errorf("Expected 1 group, got %d", resp.GroupCount)
	}
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)
	user := createTestUser(t, db, "user@test.com", "Test User", models.SystemRoleUser)

	r.PUT("/admin/users/:id", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.UpdateUser(c)
	})

	newName := "Updated Name"
	newRole := "admin"
	body, _ := json.Marshal(UpdateUserRequest{
		Name:       &newName,
		SystemRole: &newRole,
	})

	req := httptest.NewRequest("PUT", "/admin/users/"+string(rune('0'+user.ID)), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp UserResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, resp.Name)
	}
	if resp.SystemRole != newRole {
		t.Errorf("Expected role %s, got %s", newRole, resp.SystemRole)
	}
}

func TestUpdateUserCannotDemoteSelf(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)

	r.PUT("/admin/users/:id", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.UpdateUser(c)
	})

	newRole := "user"
	body, _ := json.Marshal(UpdateUserRequest{
		SystemRole: &newRole,
	})

	req := httptest.NewRequest("PUT", "/admin/users/"+string(rune('0'+admin.ID)), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)
	user := createTestUser(t, db, "user@test.com", "Test User", models.SystemRoleUser)

	// Add some related data
	group := createTestGroup(t, db, "Test Group")
	createTestLink(t, db, user.ID, group.ID, "link1")
	db.Create(&models.GroupMembership{UserID: user.ID, GroupID: group.ID, Role: models.GroupRoleMember})

	r.DELETE("/admin/users/:id", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.DeleteUser(c)
	})

	req := httptest.NewRequest("DELETE", "/admin/users/"+string(rune('0'+user.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify user is deleted
	var count int64
	db.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Errorf("Expected user to be deleted, but still exists")
	}
}

func TestDeleteUserCannotDeleteSelf(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)

	r.DELETE("/admin/users/:id", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.DeleteUser(c)
	})

	req := httptest.NewRequest("DELETE", "/admin/users/"+string(rune('0'+admin.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetStats(t *testing.T) {
	db := setupTestDB(t)
	r := setupTestRouter(db)
	h := NewHandler(db)

	admin := createTestUser(t, db, "admin@test.com", "Admin", models.SystemRoleAdmin)
	user := createTestUser(t, db, "user@test.com", "User", models.SystemRoleUser)

	group := createTestGroup(t, db, "Test Group")

	link1 := createTestLink(t, db, user.ID, group.ID, "link1")
	link1.IsPublic = true
	link1.IsUnread = false // Explicitly set to false (default is true)
	link1.ClickCount = 10
	db.Save(link1)

	link2 := createTestLink(t, db, user.ID, group.ID, "link2")
	// link2 keeps IsUnread = true (the default)
	link2.ClickCount = 5
	db.Save(link2)

	// Create a tag
	db.Create(&models.Tag{Name: "test"})

	r.GET("/admin/stats", func(c *gin.Context) {
		c.Set(auth.ContextKeyUserID, admin.ID)
		c.Set(auth.ContextKeySystemRole, "admin")
		h.GetStats(c)
	})

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var stats StatsResponse
	json.Unmarshal(w.Body.Bytes(), &stats)

	if stats.TotalUsers != 2 {
		t.Errorf("Expected 2 users, got %d", stats.TotalUsers)
	}
	if stats.TotalLinks != 2 {
		t.Errorf("Expected 2 links, got %d", stats.TotalLinks)
	}
	if stats.TotalGroups != 1 {
		t.Errorf("Expected 1 group, got %d", stats.TotalGroups)
	}
	if stats.TotalTags != 1 {
		t.Errorf("Expected 1 tag, got %d", stats.TotalTags)
	}
	if stats.TotalClicks != 15 {
		t.Errorf("Expected 15 clicks, got %d", stats.TotalClicks)
	}
	if stats.PublicLinks != 1 {
		t.Errorf("Expected 1 public link, got %d", stats.PublicLinks)
	}
	if stats.PrivateLinks != 1 {
		t.Errorf("Expected 1 private link, got %d", stats.PrivateLinks)
	}
	if stats.UnreadLinks != 1 {
		t.Errorf("Expected 1 unread link, got %d", stats.UnreadLinks)
	}
	if stats.AdminUsers != 1 {
		t.Errorf("Expected 1 admin user, got %d", stats.AdminUsers)
	}
}
