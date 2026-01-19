package apikeys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

const (
	// KeyLength is the length of the generated API key in bytes (32 bytes = 64 hex chars)
	KeyLength = 32
	// KeyPrefixLength is the number of characters to store as prefix for identification
	KeyPrefixLength = 8
)

// Handler handles API key requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new API keys handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// APIKeyResponse represents an API key in responses
type APIKeyResponse struct {
	ID          uint       `json:"id"`
	KeyPrefix   string     `json:"key_prefix"`
	Description string     `json:"description"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateAPIKeyRequest represents a request to create an API key
type CreateAPIKeyRequest struct {
	Description string `json:"description"`
}

// CreateAPIKeyResponse includes the full key (only shown once)
type CreateAPIKeyResponse struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	KeyPrefix   string    `json:"key_prefix"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// generateAPIKey generates a new random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, KeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashAPIKey creates a SHA-256 hash of the API key
func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// Create creates a new API key for the authenticated user
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Description is optional, so binding might fail with empty body
		req.Description = ""
	}

	// Generate the key
	key, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Create the API key record
	apiKey := models.APIKey{
		UserID:      userID,
		KeyHash:     hashAPIKey(key),
		KeyPrefix:   key[:KeyPrefixLength],
		Description: req.Description,
		CreatedByID: userID,
	}

	if err := h.db.Create(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	// Return the full key - this is the only time it's visible
	c.JSON(http.StatusCreated, CreateAPIKeyResponse{
		ID:          apiKey.ID,
		Key:         key,
		KeyPrefix:   apiKey.KeyPrefix,
		Description: apiKey.Description,
		CreatedAt:   apiKey.CreatedAt,
	})
}

// List returns all API keys for the authenticated user
func (h *Handler) List(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var apiKeys []models.APIKey
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&apiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch API keys"})
		return
	}

	responses := make([]APIKeyResponse, len(apiKeys))
	for i, key := range apiKeys {
		responses[i] = APIKeyResponse{
			ID:          key.ID,
			KeyPrefix:   key.KeyPrefix,
			Description: key.Description,
			LastUsedAt:  key.LastUsedAt,
			CreatedAt:   key.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// Delete deletes an API key
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	keyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	// Find the API key
	var apiKey models.APIKey
	if err := h.db.Where("id = ? AND user_id = ?", keyID, userID).First(&apiKey).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Soft delete
	if err := h.db.Delete(&apiKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted"})
}

// ValidateAPIKey checks if an API key is valid and returns the user ID
func ValidateAPIKey(db *gorm.DB, key string) (*models.APIKey, error) {
	keyHash := hashAPIKey(key)

	var apiKey models.APIKey
	if err := db.Where("key_hash = ?", keyHash).First(&apiKey).Error; err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// UpdateLastUsed updates the last_used_at timestamp for an API key
func UpdateLastUsed(db *gorm.DB, apiKeyID uint) {
	now := time.Now()
	db.Model(&models.APIKey{}).Where("id = ?", apiKeyID).Update("last_used_at", now)
}

// CombinedAuthMiddleware returns a middleware that authenticates via JWT or API key
// Both are passed in the Authorization header as "Bearer <token>"
// JWTs contain dots, API keys are hex strings without dots
func CombinedAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check for Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Try JWT first (JWTs contain dots)
		if strings.Contains(token, ".") {
			// Validate JWT
			claims, err := auth.ValidateToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.SystemRole)
			c.Next()
			return
		}

		// Try API key (hex string without dots)
		apiKey, err := ValidateAPIKey(db, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// Update last used (fire and forget)
		go UpdateLastUsed(db, apiKey.ID)

		// Set user context
		c.Set("user_id", apiKey.UserID)

		// Get user to set other context values
		var user models.User
		if err := db.First(&user, apiKey.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user_email", user.Email)
		c.Set("user_role", string(user.SystemRole))

		c.Next()
	}
}

// RegisterRoutes registers API key routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/api-keys", h.Create)
	rg.GET("/api-keys", h.List)
	rg.DELETE("/api-keys/:id", h.Delete)
}
