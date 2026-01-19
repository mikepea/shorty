package apikeys

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	handler.RegisterRoutes(api)

	return r
}

func getAuthHeader(user models.User) string {
	token, _ := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	return "Bearer " + token
}

func TestCreateAPIKey(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	body := CreateAPIKeyRequest{
		Description: "Test API Key",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/api-keys", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}

	var response CreateAPIKeyResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if response.Key == "" {
		t.Error("Expected API key to be returned")
	}

	if len(response.Key) != KeyLength*2 { // hex encoding doubles the length
		t.Errorf("Expected key length %d, got %d", KeyLength*2, len(response.Key))
	}

	if response.KeyPrefix != response.Key[:KeyPrefixLength] {
		t.Error("Key prefix should match the start of the key")
	}

	if response.Description != "Test API Key" {
		t.Errorf("Expected description 'Test API Key', got '%s'", response.Description)
	}
}

func TestCreateAPIKeyWithoutDescription(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	req, _ := http.NewRequest("POST", "/api/api-keys", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestListAPIKeys(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create a couple API keys directly in the database
	keys := []models.APIKey{
		{UserID: user.ID, KeyHash: "hash1", KeyPrefix: "key1abcd", Description: "Key 1", CreatedByID: user.ID},
		{UserID: user.ID, KeyHash: "hash2", KeyPrefix: "key2efgh", Description: "Key 2", CreatedByID: user.ID},
	}
	for _, key := range keys {
		db.Create(&key)
	}

	req, _ := http.NewRequest("GET", "/api/api-keys", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response []APIKeyResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if len(response) != 2 {
		t.Errorf("Expected 2 API keys, got %d", len(response))
	}
}

func TestListAPIKeysOnlyShowsOwnKeys(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user1 := createTestUser(t, db, "user1@example.com")
	user2 := createTestUser(t, db, "user2@example.com")

	// Create keys for both users
	db.Create(&models.APIKey{UserID: user1.ID, KeyHash: "hash1", KeyPrefix: "key1abcd", CreatedByID: user1.ID})
	db.Create(&models.APIKey{UserID: user2.ID, KeyHash: "hash2", KeyPrefix: "key2efgh", CreatedByID: user2.ID})

	// User1 should only see their own key
	req, _ := http.NewRequest("GET", "/api/api-keys", nil)
	req.Header.Set("Authorization", getAuthHeader(user1))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	var response []APIKeyResponse
	json.Unmarshal(resp.Body.Bytes(), &response)

	if len(response) != 1 {
		t.Errorf("Expected 1 API key, got %d", len(response))
	}

	if response[0].KeyPrefix != "key1abcd" {
		t.Error("Should only see own API key")
	}
}

func TestDeleteAPIKey(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user := createTestUser(t, db, "test@example.com")

	// Create an API key
	apiKey := models.APIKey{UserID: user.ID, KeyHash: "hash1", KeyPrefix: "key1abcd", CreatedByID: user.ID}
	db.Create(&apiKey)

	req, _ := http.NewRequest("DELETE", "/api/api-keys/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify it's deleted
	var count int64
	db.Model(&models.APIKey{}).Where("id = ?", apiKey.ID).Count(&count)
	if count != 0 {
		t.Error("API key should be deleted")
	}
}

func TestDeleteAPIKeyNotOwned(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestRouter(db)
	user1 := createTestUser(t, db, "user1@example.com")
	user2 := createTestUser(t, db, "user2@example.com")

	// Create an API key for user2
	apiKey := models.APIKey{UserID: user2.ID, KeyHash: "hash1", KeyPrefix: "key1abcd", CreatedByID: user2.ID}
	db.Create(&apiKey)

	// User1 tries to delete it
	req, _ := http.NewRequest("DELETE", "/api/api-keys/1", nil)
	req.Header.Set("Authorization", getAuthHeader(user1))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.Code)
	}
}

func TestValidateAPIKey(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "test@example.com")

	// Create an API key with known hash
	key := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	apiKey := models.APIKey{
		UserID:      user.ID,
		KeyHash:     hashAPIKey(key),
		KeyPrefix:   key[:KeyPrefixLength],
		CreatedByID: user.ID,
	}
	db.Create(&apiKey)

	// Validate with correct key
	result, err := ValidateAPIKey(db, key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.ID != apiKey.ID {
		t.Error("Expected to find the API key")
	}

	// Validate with wrong key
	_, err = ValidateAPIKey(db, "wrongkey")
	if err == nil {
		t.Error("Expected error for invalid key")
	}
}

func TestCombinedAuthMiddleware(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "test@example.com")

	// Create an API key
	key := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	apiKey := models.APIKey{
		UserID:      user.ID,
		KeyHash:     hashAPIKey(key),
		KeyPrefix:   key[:KeyPrefixLength],
		CreatedByID: user.ID,
	}
	db.Create(&apiKey)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CombinedAuthMiddleware(db))
	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Test with valid API key
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+key)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify user ID is set
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if uint(response["user_id"].(float64)) != user.ID {
		t.Error("User ID should be set in context")
	}
}

func TestCombinedAuthMiddlewareInvalidKey(t *testing.T) {
	db := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CombinedAuthMiddleware(db))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidkey")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.Code)
	}
}

func TestUpdateLastUsed(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, "test@example.com")

	// Create an API key without last_used
	apiKey := models.APIKey{
		UserID:      user.ID,
		KeyHash:     "hash1",
		KeyPrefix:   "key1abcd",
		CreatedByID: user.ID,
	}
	db.Create(&apiKey)

	// Update last used
	UpdateLastUsed(db, apiKey.ID)

	// Give it a moment for the update
	time.Sleep(10 * time.Millisecond)

	// Check it was updated
	var updated models.APIKey
	db.First(&updated, apiKey.ID)

	if updated.LastUsedAt == nil {
		t.Error("LastUsedAt should be set")
	}
}
