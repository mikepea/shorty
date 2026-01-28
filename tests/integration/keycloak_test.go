package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/admin"
	"github.com/mikepea/shorty/pkg/shorty/apikeys"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/groups"
	"github.com/mikepea/shorty/pkg/shorty/importexport"
	"github.com/mikepea/shorty/pkg/shorty/links"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"github.com/mikepea/shorty/pkg/shorty/oidc"
	"github.com/mikepea/shorty/pkg/shorty/redirect"
	"github.com/mikepea/shorty/pkg/shorty/scim"
	"github.com/mikepea/shorty/pkg/shorty/tags"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	keycloakContainerName = "shorty-test-keycloak"
	keycloakPort          = "8180"
	keycloakAdminUser     = "admin"
	keycloakAdminPassword = "admin"
	testRealm             = "shorty-test"
	testClientID          = "shorty"
	testClientSecret      = "shorty-secret"
)

// setupKeycloakTestDB creates an in-memory SQLite database for testing
func setupKeycloakTestDB(t *testing.T) *gorm.DB {
	// Use shared cache mode to ensure all connections share the same in-memory database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ensure single connection to prevent SQLite issues
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := models.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create the global organization if it doesn't exist (required for SCIM and redirect functionality)
	var existingOrg models.Organization
	if err := db.Where("is_global = ?", true).First(&existingOrg).Error; err != nil {
		globalOrg := models.Organization{
			Name:     "Shorty Global",
			Slug:     "shorty-global",
			IsGlobal: true,
		}
		if err := db.Create(&globalOrg).Error; err != nil {
			t.Fatalf("Failed to create global organization: %v", err)
		}
	}

	return db
}

// setupFullServerWithOIDCAndSCIM creates a Gin engine with all routes including OIDC and SCIM
func setupFullServerWithOIDCAndSCIM(db *gorm.DB, baseURL string) *gin.Engine {
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
			c.JSON(200, gin.H{"status": "ok", "service": "shorty"})
		})

		// Auth routes (public)
		authHandler := auth.NewHandler(db)
		authHandler.RegisterRoutes(api.Group("/auth"))

		// Combined auth middleware
		combinedAuth := apikeys.CombinedAuthMiddleware(db)

		// API keys routes
		apiKeysHandler := apikeys.NewHandler(db)
		apiKeysHandler.RegisterRoutes(api.Group("", auth.AuthMiddleware()))

		// Groups routes
		groupsHandler := groups.NewHandler(db)
		groupsGroup := api.Group("/groups")
		groupsGroup.Use(combinedAuth)
		groupsHandler.RegisterRoutes(groupsGroup)
		groupsHandler.RegisterMemberRoutes(groupsGroup)

		// Links routes
		linksHandler := links.NewHandler(db)
		linksHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Tags routes
		tagsHandler := tags.NewHandler(db)
		tagsHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Import/Export routes
		importExportHandler := importexport.NewHandler(db)
		importExportHandler.RegisterRoutes(api.Group("", combinedAuth))

		// Admin routes
		adminHandler := admin.NewHandler(db)
		adminGroup := api.Group("/admin")
		adminGroup.Use(auth.AuthMiddleware(), auth.RequireAdmin())
		adminHandler.RegisterRoutes(adminGroup)

		// OIDC routes
		oidcHandler := oidc.NewHandler(db, baseURL)
		oidcHandler.RegisterRoutes(api.Group("/oidc"))
		oidcHandler.RegisterAdminRoutes(adminGroup.Group("/oidc"))

		// SCIM token management (admin only)
		scimTokenHandler := scim.NewTokenHandler(db)
		scimTokenHandler.RegisterAdminRoutes(adminGroup)
	}

	// SCIM routes
	scimGroup := r.Group("/scim/v2")
	scimGroup.Use(scim.SCIMAuthMiddleware(db))
	{
		scimUserHandler := scim.NewUserHandler(db, baseURL)
		scimUserHandler.RegisterRoutes(scimGroup)

		scimGroupHandler := scim.NewGroupHandler(db, baseURL)
		scimGroupHandler.RegisterRoutes(scimGroup)

		scimConfigHandler := scim.NewConfigHandler(db, baseURL)
		scimConfigHandler.RegisterRoutes(scimGroup)
	}

	// Redirect routes
	redirectHandler := redirect.NewHandler(db)
	redirectHandler.RegisterRoutes(r)

	return r
}

// createAdminUser creates an admin user and returns a JWT token
func createAdminUser(t *testing.T, db *gorm.DB, router *gin.Engine) string {
	// Create admin user
	hashedPassword, _ := auth.HashPassword("adminpass")
	adminUser := models.User{
		Email:        "admin@test.com",
		Name:         "Admin",
		PasswordHash: hashedPassword,
		SystemRole:   models.SystemRoleAdmin,
		Active:       true,
	}
	db.Create(&adminUser)

	// Login to get token
	loginBody := `{"email":"admin@test.com","password":"adminpass"}`
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Failed to login as admin: %d - %s", resp.Code, resp.Body.String())
	}

	var authResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(resp.Body.Bytes(), &authResp)
	return authResp.Token
}

// isDockerAvailable checks if Docker is available
func isDockerAvailable() bool {
	cmd := exec.Command("docker", "version")
	return cmd.Run() == nil
}

// startKeycloak starts Keycloak in a Docker container
func startKeycloak(t *testing.T) (cleanup func()) {
	// Check if Keycloak is already running
	checkCmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", keycloakContainerName)
	output, err := checkCmd.Output()
	if err == nil && strings.TrimSpace(string(output)) == "true" {
		t.Log("Keycloak already running, reusing existing container")
		return func() {} // No cleanup needed
	}

	// Remove any existing stopped container
	exec.Command("docker", "rm", "-f", keycloakContainerName).Run()

	// Start Keycloak
	t.Log("Starting Keycloak container...")
	cmd := exec.Command("docker", "run", "-d",
		"--name", keycloakContainerName,
		"-p", keycloakPort+":8080",
		"-e", "KEYCLOAK_ADMIN="+keycloakAdminUser,
		"-e", "KEYCLOAK_ADMIN_PASSWORD="+keycloakAdminPassword,
		"quay.io/keycloak/keycloak:latest",
		"start-dev",
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start Keycloak: %v", err)
	}

	// Wait for Keycloak to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	t.Log("Waiting for Keycloak to be ready...")
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Keycloak did not become ready in time")
		default:
			resp, err := http.Get("http://localhost:" + keycloakPort + "/health/ready")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				t.Log("Keycloak is ready")
				return func() {
					exec.Command("docker", "rm", "-f", keycloakContainerName).Run()
				}
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(2 * time.Second)
		}
	}
}

// getKeycloakAdminToken gets an admin token from Keycloak
func getKeycloakAdminToken(t *testing.T) string {
	resp, err := http.PostForm("http://localhost:"+keycloakPort+"/realms/master/protocol/openid-connect/token",
		map[string][]string{
			"client_id":  {"admin-cli"},
			"username":   {keycloakAdminUser},
			"password":   {keycloakAdminPassword},
			"grant_type": {"password"},
		})
	if err != nil {
		t.Fatalf("Failed to get admin token: %v", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	return tokenResp.AccessToken
}

// setupKeycloakRealm creates a test realm with a client
func setupKeycloakRealm(t *testing.T, adminToken string, callbackURL string) {
	client := &http.Client{Timeout: 30 * time.Second}

	// Create realm
	realmData := map[string]interface{}{
		"realm":   testRealm,
		"enabled": true,
	}
	realmJSON, _ := json.Marshal(realmData)

	req, _ := http.NewRequest("POST", "http://localhost:"+keycloakPort+"/admin/realms", bytes.NewReader(realmJSON))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create realm: %v", err)
	}
	resp.Body.Close()
	// 409 is OK - realm already exists
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		t.Fatalf("Failed to create realm: %d", resp.StatusCode)
	}

	// Create client
	clientData := map[string]interface{}{
		"clientId":                  testClientID,
		"secret":                    testClientSecret,
		"enabled":                   true,
		"directAccessGrantsEnabled": true,
		"standardFlowEnabled":       true,
		"redirectUris":              []string{callbackURL},
		"webOrigins":                []string{"*"},
	}
	clientJSON, _ := json.Marshal(clientData)

	req, _ = http.NewRequest("POST", "http://localhost:"+keycloakPort+"/admin/realms/"+testRealm+"/clients", bytes.NewReader(clientJSON))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	resp.Body.Close()

	// Create test user
	userData := map[string]interface{}{
		"username":      "testuser",
		"email":         "testuser@example.com",
		"emailVerified": true,
		"enabled":       true,
		"firstName":     "Test",
		"lastName":      "User",
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     "testpass",
				"temporary": false,
			},
		},
	}
	userJSON, _ := json.Marshal(userData)

	req, _ = http.NewRequest("POST", "http://localhost:"+keycloakPort+"/admin/realms/"+testRealm+"/users", bytes.NewReader(userJSON))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	resp.Body.Close()

	t.Log("Keycloak realm, client, and user created")
}

// TestKeycloakOIDCIntegration tests OIDC with Keycloak
func TestKeycloakOIDCIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST_KEYCLOAK") != "1" {
		t.Skip("Skipping Keycloak integration test. Set INTEGRATION_TEST_KEYCLOAK=1 to run.")
	}

	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping Keycloak test")
	}

	cleanup := startKeycloak(t)
	defer cleanup()

	db := setupKeycloakTestDB(t)
	baseURL := "http://localhost:8080"
	router := setupFullServerWithOIDCAndSCIM(db, baseURL)

	// Create admin and get token
	adminToken := createAdminUser(t, db, router)

	// Get Keycloak admin token and setup realm
	keycloakToken := getKeycloakAdminToken(t)
	setupKeycloakRealm(t, keycloakToken, baseURL+"/api/oidc/callback")

	// Create OIDC provider in Shorty
	t.Run("CreateOIDCProvider", func(t *testing.T) {
		providerData := map[string]interface{}{
			"name":           "Keycloak",
			"slug":           "keycloak",
			"issuer":         "http://localhost:" + keycloakPort + "/realms/" + testRealm,
			"client_id":      testClientID,
			"client_secret":  testClientSecret,
			"scopes":         "openid profile email",
			"enabled":        true,
			"auto_provision": true,
		}
		body, _ := json.Marshal(providerData)

		req, _ := http.NewRequest("POST", "/api/admin/oidc/providers", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d: %s", resp.Code, resp.Body.String())
		}
	})

	// List OIDC providers
	t.Run("ListOIDCProviders", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/oidc/providers", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.Code)
		}

		var providers []map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &providers)

		if len(providers) != 1 {
			t.Errorf("Expected 1 provider, got %d", len(providers))
		}
	})

	// Get auth URL
	t.Run("GetAuthURL", func(t *testing.T) {
		body := `{"return_url":"http://localhost:3000/sso/callback"}`
		req, _ := http.NewRequest("POST", "/api/oidc/providers/keycloak/auth", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", resp.Code, resp.Body.String())
		}

		var authResp map[string]string
		json.Unmarshal(resp.Body.Bytes(), &authResp)

		if !strings.Contains(authResp["auth_url"], "localhost:"+keycloakPort) {
			t.Errorf("Auth URL should contain Keycloak URL: %s", authResp["auth_url"])
		}
	})
}

// TestSCIMIntegration tests SCIM endpoints
func TestSCIMIntegration(t *testing.T) {
	db := setupKeycloakTestDB(t)
	baseURL := "http://localhost:8080"
	router := setupFullServerWithOIDCAndSCIM(db, baseURL)

	// Create admin and get token
	adminToken := createAdminUser(t, db, router)

	var scimToken string

	// Create SCIM token (organization_id 1 is the global org created in setup)
	t.Run("CreateSCIMToken", func(t *testing.T) {
		body := `{"organization_id":1,"description":"Integration Test Token"}`
		req, _ := http.NewRequest("POST", "/api/admin/scim-tokens", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected 201, got %d: %s", resp.Code, resp.Body.String())
		}

		var tokenResp struct {
			Token string `json:"token"`
		}
		json.Unmarshal(resp.Body.Bytes(), &tokenResp)
		scimToken = tokenResp.Token
	})

	// Test ServiceProviderConfig
	t.Run("ServiceProviderConfig", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/scim/v2/ServiceProviderConfig", nil)
		req.Header.Set("Authorization", "Bearer "+scimToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.Code)
		}

		var config map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &config)

		if config["patch"] == nil {
			t.Error("Expected patch configuration")
		}
	})

	// Test SCIM User lifecycle
	var createdUserID string

	t.Run("CreateSCIMUser", func(t *testing.T) {
		userData := map[string]interface{}{
			"schemas":     []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			"userName":    "scim.user@example.com",
			"displayName": "SCIM User",
			"externalId":  "ext-user-123",
			"name": map[string]string{
				"givenName":  "SCIM",
				"familyName": "User",
			},
			"emails": []map[string]interface{}{
				{"value": "scim.user@example.com", "primary": true},
			},
			"active": true,
		}
		body, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/scim/v2/Users", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+scimToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected 201, got %d: %s", resp.Code, resp.Body.String())
		}

		var user map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &user)
		createdUserID = user["id"].(string)
	})

	t.Run("GetSCIMUser", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/scim/v2/Users/"+createdUserID, nil)
		req.Header.Set("Authorization", "Bearer "+scimToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.Code)
		}

		var user map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &user)

		if user["userName"] != "scim.user@example.com" {
			t.Errorf("Expected userName scim.user@example.com, got %s", user["userName"])
		}
	})

	t.Run("FilterSCIMUsers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/scim/v2/Users?filter=userName%20eq%20%22scim.user@example.com%22", nil)
		req.Header.Set("Authorization", "Bearer "+scimToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.Code)
		}

		var listResp map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &listResp)

		totalResults := int(listResp["totalResults"].(float64))
		if totalResults != 1 {
			t.Errorf("Expected 1 result, got %d", totalResults)
		}
	})

	t.Run("PatchSCIMUserDeactivate", func(t *testing.T) {
		patchData := map[string]interface{}{
			"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
			"Operations": []map[string]interface{}{
				{
					"op":    "replace",
					"path":  "active",
					"value": false,
				},
			},
		}
		body, _ := json.Marshal(patchData)

		req, _ := http.NewRequest("PATCH", "/scim/v2/Users/"+createdUserID, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+scimToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", resp.Code, resp.Body.String())
		}

		var user map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &user)

		if user["active"] != false {
			t.Error("Expected user to be deactivated")
		}
	})

	// Test SCIM Group lifecycle
	var createdGroupID string

	t.Run("CreateSCIMGroup", func(t *testing.T) {
		groupData := map[string]interface{}{
			"schemas":     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
			"displayName": "SCIM Test Group",
			"externalId":  "ext-group-123",
		}
		body, _ := json.Marshal(groupData)

		req, _ := http.NewRequest("POST", "/scim/v2/Groups", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+scimToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("Expected 201, got %d: %s", resp.Code, resp.Body.String())
		}

		var group map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &group)
		createdGroupID = group["id"].(string)
	})

	t.Run("PatchSCIMGroupAddMember", func(t *testing.T) {
		if createdGroupID == "" {
			t.Skip("Skipping: group was not created")
		}

		patchData := map[string]interface{}{
			"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
			"Operations": []map[string]interface{}{
				{
					"op":   "add",
					"path": "members",
					"value": []map[string]string{
						{"value": createdUserID},
					},
				},
			},
		}
		body, _ := json.Marshal(patchData)

		req, _ := http.NewRequest("PATCH", "/scim/v2/Groups/"+createdGroupID, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+scimToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", resp.Code, resp.Body.String())
			return
		}

		var group map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &group)

		members, ok := group["members"].([]interface{})
		if !ok || len(members) != 1 {
			t.Errorf("Expected 1 member, got %d", len(members))
		}
	})

	t.Run("PatchSCIMGroupRemoveMember", func(t *testing.T) {
		if createdGroupID == "" {
			t.Skip("Skipping: group was not created")
		}

		patchData := map[string]interface{}{
			"schemas": []string{"urn:ietf:params:scim:api:messages:2.0:PatchOp"},
			"Operations": []map[string]interface{}{
				{
					"op":   "remove",
					"path": fmt.Sprintf("members[value eq \"%s\"]", createdUserID),
				},
			},
		}
		body, _ := json.Marshal(patchData)

		req, _ := http.NewRequest("PATCH", "/scim/v2/Groups/"+createdGroupID, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+scimToken)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d: %s", resp.Code, resp.Body.String())
			return
		}

		var group map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &group)

		members, ok := group["members"].([]interface{})
		if ok && len(members) != 0 {
			t.Errorf("Expected 0 members, got %d", len(members))
		}
	})

	t.Run("DeleteSCIMGroup", func(t *testing.T) {
		if createdGroupID == "" {
			t.Skip("Skipping: group was not created")
		}

		req, _ := http.NewRequest("DELETE", "/scim/v2/Groups/"+createdGroupID, nil)
		req.Header.Set("Authorization", "Bearer "+scimToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", resp.Code)
		}
	})

	t.Run("DeleteSCIMUser", func(t *testing.T) {
		if createdUserID == "" {
			t.Skip("Skipping: user was not created")
		}

		req, _ := http.NewRequest("DELETE", "/scim/v2/Users/"+createdUserID, nil)
		req.Header.Set("Authorization", "Bearer "+scimToken)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", resp.Code)
		}
	})
}

// TestSCIMAuthRequired tests that SCIM endpoints require authentication
func TestSCIMAuthRequired(t *testing.T) {
	db := setupKeycloakTestDB(t)
	baseURL := "http://localhost:8080"
	router := setupFullServerWithOIDCAndSCIM(db, baseURL)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/scim/v2/Users"},
		{"POST", "/scim/v2/Users"},
		{"GET", "/scim/v2/Users/1"},
		{"PUT", "/scim/v2/Users/1"},
		{"PATCH", "/scim/v2/Users/1"},
		{"DELETE", "/scim/v2/Users/1"},
		{"GET", "/scim/v2/Groups"},
		{"POST", "/scim/v2/Groups"},
		{"GET", "/scim/v2/ServiceProviderConfig"},
	}

	for _, e := range endpoints {
		t.Run(e.method+" "+e.path, func(t *testing.T) {
			req, _ := http.NewRequest(e.method, e.path, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("Expected 401 for %s %s without auth, got %d", e.method, e.path, resp.Code)
			}
		})
	}
}
