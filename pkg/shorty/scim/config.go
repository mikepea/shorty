package scim

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// ConfigHandler handles SCIM configuration endpoints
type ConfigHandler struct {
	db      *gorm.DB
	baseURL string
}

// NewConfigHandler creates a new SCIM config handler
func NewConfigHandler(db *gorm.DB, baseURL string) *ConfigHandler {
	return &ConfigHandler{db: db, baseURL: baseURL}
}

// GetServiceProviderConfig returns the service provider configuration
func (h *ConfigHandler) GetServiceProviderConfig(c *gin.Context) {
	c.JSON(http.StatusOK, ServiceProviderConfig{
		Schemas: []string{SchemaServiceProvider},
		Patch: SupportedConfig{
			Supported: true,
		},
		Bulk: BulkConfig{
			Supported:      false,
			MaxOperations:  0,
			MaxPayloadSize: 0,
		},
		Filter: FilterConfig{
			Supported:  true,
			MaxResults: 1000,
		},
		ChangePassword: SupportedConfig{
			Supported: false,
		},
		Sort: SupportedConfig{
			Supported: false,
		},
		Etag: SupportedConfig{
			Supported: false,
		},
		AuthenticationSchemes: []AuthenticationScheme{
			{
				Type:        "oauthbearertoken",
				Name:        "OAuth Bearer Token",
				Description: "Authentication scheme using the OAuth Bearer Token Standard",
				SpecURI:     "https://www.rfc-editor.org/info/rfc6750",
				Primary:     true,
			},
		},
		Meta: Meta{
			ResourceType: "ServiceProviderConfig",
			Location:     h.baseURL + "/scim/v2/ServiceProviderConfig",
		},
	})
}

// GetResourceTypes returns the supported resource types
func (h *ConfigHandler) GetResourceTypes(c *gin.Context) {
	c.JSON(http.StatusOK, []ResourceType{
		{
			Schemas:     []string{SchemaResourceType},
			ID:          "User",
			Name:        "User",
			Endpoint:    "/Users",
			Description: "User Account",
			Schema:      SchemaUser,
			Meta: Meta{
				ResourceType: "ResourceType",
				Location:     h.baseURL + "/scim/v2/ResourceTypes/User",
			},
		},
		{
			Schemas:     []string{SchemaResourceType},
			ID:          "Group",
			Name:        "Group",
			Endpoint:    "/Groups",
			Description: "Group",
			Schema:      SchemaGroup,
			Meta: Meta{
				ResourceType: "ResourceType",
				Location:     h.baseURL + "/scim/v2/ResourceTypes/Group",
			},
		},
	})
}

// GetSchemas returns the supported schemas
func (h *ConfigHandler) GetSchemas(c *gin.Context) {
	// Return a simplified schema response
	c.JSON(http.StatusOK, ListResponse{
		Schemas:      []string{SchemaListResponse},
		TotalResults: 2,
		StartIndex:   1,
		ItemsPerPage: 2,
		Resources: []map[string]interface{}{
			{
				"id":          SchemaUser,
				"name":        "User",
				"description": "User Account",
				"attributes":  getUserSchemaAttributes(),
				"meta": map[string]string{
					"resourceType": "Schema",
					"location":     h.baseURL + "/scim/v2/Schemas/" + SchemaUser,
				},
			},
			{
				"id":          SchemaGroup,
				"name":        "Group",
				"description": "Group",
				"attributes":  getGroupSchemaAttributes(),
				"meta": map[string]string{
					"resourceType": "Schema",
					"location":     h.baseURL + "/scim/v2/Schemas/" + SchemaGroup,
				},
			},
		},
	})
}

func getUserSchemaAttributes() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "userName", "type": "string", "multiValued": false, "required": true, "caseExact": false, "mutability": "readWrite", "returned": "default", "uniqueness": "server"},
		{"name": "name", "type": "complex", "multiValued": false, "required": false, "mutability": "readWrite", "returned": "default"},
		{"name": "displayName", "type": "string", "multiValued": false, "required": false, "mutability": "readWrite", "returned": "default"},
		{"name": "emails", "type": "complex", "multiValued": true, "required": false, "mutability": "readWrite", "returned": "default"},
		{"name": "active", "type": "boolean", "multiValued": false, "required": false, "mutability": "readWrite", "returned": "default"},
		{"name": "externalId", "type": "string", "multiValued": false, "required": false, "caseExact": true, "mutability": "readWrite", "returned": "default"},
	}
}

func getGroupSchemaAttributes() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "displayName", "type": "string", "multiValued": false, "required": true, "mutability": "readWrite", "returned": "default"},
		{"name": "members", "type": "complex", "multiValued": true, "required": false, "mutability": "readWrite", "returned": "default"},
		{"name": "externalId", "type": "string", "multiValued": false, "required": false, "caseExact": true, "mutability": "readWrite", "returned": "default"},
	}
}

// RegisterRoutes registers SCIM config routes
func (h *ConfigHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ServiceProviderConfig", h.GetServiceProviderConfig)
	rg.GET("/ResourceTypes", h.GetResourceTypes)
	rg.GET("/Schemas", h.GetSchemas)
}

// SCIM Token management

// GenerateSCIMToken creates a new SCIM bearer token for an organization
func GenerateSCIMToken(db *gorm.DB, organizationID uint, description string) (string, *models.SCIMToken, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash the token
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	scimToken := &models.SCIMToken{
		OrganizationID: organizationID,
		TokenHash:      tokenHash,
		TokenPrefix:    token[:8],
		Description:    description,
	}

	if err := db.Create(scimToken).Error; err != nil {
		return "", nil, err
	}

	return token, scimToken, nil
}

// ValidateSCIMToken validates a SCIM bearer token
func ValidateSCIMToken(db *gorm.DB, token string) (*models.SCIMToken, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var scimToken models.SCIMToken
	if err := db.Where("token_hash = ?", tokenHash).First(&scimToken).Error; err != nil {
		return nil, err
	}

	// Update last used (fire and forget)
	go func() {
		now := time.Now()
		db.Model(&scimToken).Update("last_used_at", &now)
	}()

	return &scimToken, nil
}

// Context key for SCIM organization ID
const ContextKeySCIMOrgID = "scim_organization_id"

// SCIMAuthMiddleware authenticates SCIM requests using bearer tokens
// and sets the organization context from the token
func SCIMAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Schemas: []string{SchemaError},
				Detail:  "Authorization header required",
				Status:  "401",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Schemas: []string{SchemaError},
				Detail:  "Invalid authorization header format",
				Status:  "401",
			})
			c.Abort()
			return
		}

		token := parts[1]
		scimToken, err := ValidateSCIMToken(db, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Schemas: []string{SchemaError},
				Detail:  "Invalid token",
				Status:  "401",
			})
			c.Abort()
			return
		}

		// Set organization context from token
		c.Set(ContextKeySCIMOrgID, scimToken.OrganizationID)

		c.Next()
	}
}

// GetSCIMOrgID returns the organization ID from SCIM context
func GetSCIMOrgID(c *gin.Context) (uint, bool) {
	orgID, exists := c.Get(ContextKeySCIMOrgID)
	if !exists {
		return 0, false
	}
	return orgID.(uint), true
}

// TokenResponse represents a SCIM token in API responses
type TokenResponse struct {
	ID          uint       `json:"id"`
	TokenPrefix string     `json:"token_prefix"`
	Description string     `json:"description"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateTokenResponse includes the full token (only shown on creation)
type CreateTokenResponse struct {
	ID          uint      `json:"id"`
	Token       string    `json:"token"`
	TokenPrefix string    `json:"token_prefix"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// TokenHandler handles SCIM token management (admin only)
type TokenHandler struct {
	db *gorm.DB
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(db *gorm.DB) *TokenHandler {
	return &TokenHandler{db: db}
}

// ListTokens returns all SCIM tokens
func (h *TokenHandler) ListTokens(c *gin.Context) {
	var tokens []models.SCIMToken
	h.db.Find(&tokens)

	responses := make([]TokenResponse, len(tokens))
	for i, t := range tokens {
		responses[i] = TokenResponse{
			ID:          t.ID,
			TokenPrefix: t.TokenPrefix,
			Description: t.Description,
			LastUsedAt:  t.LastUsedAt,
			CreatedAt:   t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// CreateTokenRequest represents a request to create a SCIM token with organization
type CreateTokenRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required"`
	Description    string `json:"description"`
}

// CreateToken creates a new SCIM token for an organization
func (h *TokenHandler) CreateToken(c *gin.Context) {
	var req CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify organization exists
	var org models.Organization
	if err := h.db.First(&org, req.OrganizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	token, scimToken, err := GenerateSCIMToken(h.db, req.OrganizationID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, CreateTokenResponse{
		ID:          scimToken.ID,
		Token:       token,
		TokenPrefix: scimToken.TokenPrefix,
		Description: scimToken.Description,
		CreatedAt:   scimToken.CreatedAt,
	})
}

// DeleteToken deletes a SCIM token
func (h *TokenHandler) DeleteToken(c *gin.Context) {
	id := c.Param("id")

	var token models.SCIMToken
	if err := h.db.First(&token, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	h.db.Delete(&token)
	c.JSON(http.StatusOK, gin.H{"message": "Token deleted"})
}

// RegisterAdminRoutes registers SCIM token admin routes
func (h *TokenHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/scim-tokens", h.ListTokens)
	rg.POST("/scim-tokens", h.CreateToken)
	rg.DELETE("/scim-tokens/:id", h.DeleteToken)
}
