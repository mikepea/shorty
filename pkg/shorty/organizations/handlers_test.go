package organizations

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

func createGlobalOrg(t *testing.T, db *gorm.DB) models.Organization {
	globalOrg := models.Organization{
		Name:     "Shorty Global",
		Slug:     "shorty-global",
		IsGlobal: true,
	}
	if err := db.Create(&globalOrg).Error; err != nil {
		t.Fatalf("Failed to create global organization: %v", err)
	}
	return globalOrg
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(db)

	orgs := r.Group("/organizations")
	orgs.Use(auth.AuthMiddleware())
	handler.RegisterRoutes(orgs)
	handler.RegisterMemberRoutes(orgs)

	return r
}

func getAuthHeader(user models.User) string {
	token, _ := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	return "Bearer " + token
}

func TestCreateOrganization(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	body := CreateOrgRequest{
		Name: "Acme Corp",
		Slug: "acme",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response OrgResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Acme Corp" {
		t.Errorf("Expected name 'Acme Corp', got %s", response.Name)
	}
	if response.Slug != "acme" {
		t.Errorf("Expected slug 'acme', got %s", response.Slug)
	}
	if response.Role != "admin" {
		t.Errorf("Expected role 'admin', got %s", response.Role)
	}
}

func TestListOrganizations(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org and add user
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/organizations", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var orgs []OrgResponse
	json.Unmarshal(resp.Body.Bytes(), &orgs)

	if len(orgs) != 1 {
		t.Errorf("Expected 1 organization, got %d", len(orgs))
	}
}

func TestGetOrganization(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org and add user
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/organizations/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var response OrgResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Test Org" {
		t.Errorf("Expected name 'Test Org', got %s", response.Name)
	}
}

func TestGetOrganizationNotMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org without adding user
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)

	req, _ := http.NewRequest("GET", "/organizations/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestUpdateOrganization(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org and add user as admin
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	body := UpdateOrgRequest{Name: "Updated Org"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/organizations/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response OrgResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Name != "Updated Org" {
		t.Errorf("Expected name 'Updated Org', got %s", response.Name)
	}
}

func TestDeleteOrganization(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org and add user as admin
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("DELETE", "/organizations/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

func TestCannotDeleteGlobalOrganization(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")
	globalOrg := createGlobalOrg(t, db)

	// Add user as admin of global org
	membership := models.OrganizationMembership{
		OrganizationID: globalOrg.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("DELETE", "/organizations/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestListMembers(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create org and add user
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	req, _ := http.NewRequest("GET", "/organizations/1/members", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
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
	newUser := createTestUser(t, db, "newuser@example.com")

	// Create org and add admin
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	membership := models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         admin.ID,
		Role:           models.OrgRoleAdmin,
	}
	db.Create(&membership)

	body := AddMemberRequest{
		Email: newUser.Email,
		Role:  "member",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/organizations/1/members", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestRemoveMember(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	admin := createTestUser(t, db, "admin@example.com")
	member := createTestUser(t, db, "member@example.com")

	// Create org and add both users
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	db.Create(&models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         admin.ID,
		Role:           models.OrgRoleAdmin,
	})
	db.Create(&models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         member.ID,
		Role:           models.OrgRoleMember,
	})

	req, _ := http.NewRequest("DELETE", "/organizations/1/members/2", nil)
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

func TestCannotRemoveOnlyAdmin(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	admin := createTestUser(t, db, "admin@example.com")

	// Create org with only one admin
	org := models.Organization{Name: "Test Org", Slug: "test-org"}
	db.Create(&org)
	db.Create(&models.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         admin.ID,
		Role:           models.OrgRoleAdmin,
	})

	// Try to remove self as only admin
	req, _ := http.NewRequest("DELETE", "/organizations/1/members/1", nil)
	req.Header.Set("Authorization", getAuthHeader(admin))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestSlugValidation(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	testCases := []struct {
		slug         string
		expectStatus int
		description  string
	}{
		{"valid-slug", http.StatusCreated, "valid slug with hyphen"},
		{"Valid123", http.StatusCreated, "uppercase normalized to lowercase"},
		{"-invalid", http.StatusBadRequest, "leading hyphen not allowed"},
		{"invalid-", http.StatusBadRequest, "trailing hyphen not allowed"},
		{"api", http.StatusBadRequest, "reserved slug"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			body := CreateOrgRequest{
				Name: "Test Org",
				Slug: tc.slug,
			}
			jsonBody, _ := json.Marshal(body)

			req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", getAuthHeader(user))
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tc.expectStatus {
				t.Errorf("Slug '%s': expected status %d, got %d: %s", tc.slug, tc.expectStatus, resp.Code, resp.Body.String())
			}
		})
	}
}
