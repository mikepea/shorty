package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// Handler handles OIDC-related requests
type Handler struct {
	db        *gorm.DB
	baseURL   string
	providers map[uint]*providerConfig
	mu        sync.RWMutex
}

type providerConfig struct {
	provider *oidc.Provider
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
}

// StateData stores OIDC state for validation
type StateData struct {
	ProviderID uint   `json:"provider_id"`
	ReturnURL  string `json:"return_url"`
	Nonce      string `json:"nonce"`
}

// NewHandler creates a new OIDC handler
func NewHandler(db *gorm.DB, baseURL string) *Handler {
	h := &Handler{
		db:        db,
		baseURL:   baseURL,
		providers: make(map[uint]*providerConfig),
	}
	// Load existing providers
	h.loadProviders()
	return h
}

// loadProviders loads all enabled OIDC providers from the database
func (h *Handler) loadProviders() {
	var providers []models.OIDCProvider
	h.db.Where("enabled = ?", true).Find(&providers)

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, p := range providers {
		if err := h.initProvider(p); err != nil {
			// Log error but continue with other providers
			continue
		}
	}
}

// initProvider initializes an OIDC provider
func (h *Handler) initProvider(p models.OIDCProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, p.Issuer)
	if err != nil {
		return err
	}

	scopes := strings.Fields(p.Scopes)
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	config := oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  h.baseURL + "/api/oidc/callback",
		Scopes:       scopes,
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: p.ClientID})

	h.providers[p.ID] = &providerConfig{
		provider: provider,
		config:   config,
		verifier: verifier,
	}

	return nil
}

// ProviderResponse represents an OIDC provider in API responses
type ProviderResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Enabled bool   `json:"enabled"`
}

// ListProviders returns all enabled OIDC providers (public endpoint)
func (h *Handler) ListProviders(c *gin.Context) {
	var providers []models.OIDCProvider
	h.db.Where("enabled = ?", true).Find(&providers)

	responses := make([]ProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = ProviderResponse{
			ID:      p.ID,
			Name:    p.Name,
			Slug:    p.Slug,
			Enabled: p.Enabled,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// AuthURLRequest represents a request for an auth URL
type AuthURLRequest struct {
	ReturnURL string `json:"return_url"`
}

// GetAuthURL returns the authorization URL for an OIDC provider
func (h *Handler) GetAuthURL(c *gin.Context) {
	slug := c.Param("slug")

	var provider models.OIDCProvider
	if err := h.db.Where("slug = ? AND enabled = ?", slug, true).First(&provider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	h.mu.RLock()
	pc, ok := h.providers[provider.ID]
	h.mu.RUnlock()

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Provider not configured"})
		return
	}

	var req AuthURLRequest
	c.ShouldBindJSON(&req)

	// Generate state with provider ID and return URL
	nonce := generateRandomString(32)
	stateData := StateData{
		ProviderID: provider.ID,
		ReturnURL:  req.ReturnURL,
		Nonce:      nonce,
	}
	stateJSON, _ := json.Marshal(stateData)
	state := base64.URLEncoding.EncodeToString(stateJSON)

	authURL := pc.config.AuthCodeURL(state, oidc.Nonce(nonce))

	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

// Callback handles the OIDC callback
func (h *Handler) Callback(c *gin.Context) {
	// Parse state
	stateParam := c.Query("state")
	stateJSON, err := base64.URLEncoding.DecodeString(stateParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	var stateData StateData
	if err := json.Unmarshal(stateJSON, &stateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Get provider config
	h.mu.RLock()
	pc, ok := h.providers[stateData.ProviderID]
	h.mu.RUnlock()

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown provider"})
		return
	}

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		errorDesc := c.Query("error_description")
		if errorDesc == "" {
			errorDesc = c.Query("error")
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failed: " + errorDesc})
		return
	}

	ctx := context.Background()
	oauth2Token, err := pc.config.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Extract ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No ID token in response"})
		return
	}

	// Verify ID token
	idToken, err := pc.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token"})
		return
	}

	// Verify nonce
	if idToken.Nonce != stateData.Nonce {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nonce"})
		return
	}

	// Extract claims
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}

	if claims.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not provided by identity provider"})
		return
	}

	// Get provider details
	var provider models.OIDCProvider
	h.db.First(&provider, stateData.ProviderID)

	// Find or create user
	user, err := h.findOrCreateUser(idToken.Subject, claims.Email, claims.Name, claims.GivenName, claims.FamilyName, &provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user: " + err.Error()})
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "User account is deactivated"})
		return
	}

	// Generate JWT
	token, err := auth.GenerateToken(user.ID, user.Email, string(user.SystemRole))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect with token or return JSON based on return URL
	if stateData.ReturnURL != "" {
		// Redirect to frontend with token
		redirectURL := stateData.ReturnURL + "?token=" + token
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": auth.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			SystemRole: string(user.SystemRole),
		},
	})
}

// findOrCreateUser finds an existing user or creates a new one
func (h *Handler) findOrCreateUser(subject, email, name, givenName, familyName string, provider *models.OIDCProvider) (*models.User, error) {
	// First, check if we have an OIDC identity link
	var identity models.OIDCIdentity
	err := h.db.Where("provider_id = ? AND subject = ?", provider.ID, subject).First(&identity).Error

	if err == nil {
		// Found existing identity, get the user
		var user models.User
		if err := h.db.First(&user, identity.UserID).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	// No identity link, check if user exists by email
	var user models.User
	err = h.db.Where("email = ?", email).First(&user).Error

	if err == nil {
		// User exists, create identity link
		identity := models.OIDCIdentity{
			UserID:     user.ID,
			ProviderID: provider.ID,
			Subject:    subject,
			Email:      email,
		}
		h.db.Create(&identity)
		return &user, nil
	}

	// User doesn't exist, check if auto-provisioning is enabled
	if !provider.AutoProvision {
		return nil, err
	}

	// Create new user
	if name == "" {
		if givenName != "" || familyName != "" {
			name = strings.TrimSpace(givenName + " " + familyName)
		} else {
			name = strings.Split(email, "@")[0]
		}
	}

	user = models.User{
		Email:      email,
		Name:       name,
		GivenName:  givenName,
		FamilyName: familyName,
		Active:     true,
		SystemRole: models.SystemRoleUser,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create OIDC identity link
		identity := models.OIDCIdentity{
			UserID:     user.ID,
			ProviderID: provider.ID,
			Subject:    subject,
			Email:      email,
		}
		if err := tx.Create(&identity).Error; err != nil {
			return err
		}

		// Create personal group
		personalGroup := models.Group{
			Name:        user.Name + "'s Links",
			Description: "Personal links for " + user.Name,
		}
		if err := tx.Create(&personalGroup).Error; err != nil {
			return err
		}

		// Add user as admin of personal group
		membership := models.GroupMembership{
			UserID:  user.ID,
			GroupID: personalGroup.ID,
			Role:    models.GroupRoleAdmin,
		}
		return tx.Create(&membership).Error
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Admin endpoints for managing OIDC providers

// AdminProviderResponse includes all provider details for admins
type AdminProviderResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Issuer        string `json:"issuer"`
	ClientID      string `json:"client_id"`
	Scopes        string `json:"scopes"`
	Enabled       bool   `json:"enabled"`
	AutoProvision bool   `json:"auto_provision"`
	CreatedAt     string `json:"created_at"`
}

// CreateProviderRequest represents a request to create an OIDC provider
type CreateProviderRequest struct {
	Name          string `json:"name" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	Issuer        string `json:"issuer" binding:"required,url"`
	ClientID      string `json:"client_id" binding:"required"`
	ClientSecret  string `json:"client_secret" binding:"required"`
	Scopes        string `json:"scopes"`
	Enabled       bool   `json:"enabled"`
	AutoProvision bool   `json:"auto_provision"`
}

// ListProvidersAdmin returns all OIDC providers for admin
func (h *Handler) ListProvidersAdmin(c *gin.Context) {
	var providers []models.OIDCProvider
	h.db.Find(&providers)

	responses := make([]AdminProviderResponse, len(providers))
	for i, p := range providers {
		responses[i] = AdminProviderResponse{
			ID:            p.ID,
			Name:          p.Name,
			Slug:          p.Slug,
			Issuer:        p.Issuer,
			ClientID:      p.ClientID,
			Scopes:        p.Scopes,
			Enabled:       p.Enabled,
			AutoProvision: p.AutoProvision,
			CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, responses)
}

// CreateProvider creates a new OIDC provider
func (h *Handler) CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := models.OIDCProvider{
		Name:          req.Name,
		Slug:          req.Slug,
		Issuer:        req.Issuer,
		ClientID:      req.ClientID,
		ClientSecret:  req.ClientSecret,
		Scopes:        req.Scopes,
		Enabled:       req.Enabled,
		AutoProvision: req.AutoProvision,
	}

	if err := h.db.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider"})
		return
	}

	// Initialize provider if enabled
	if provider.Enabled {
		if err := h.initProvider(provider); err != nil {
			// Provider created but initialization failed - log warning
			c.JSON(http.StatusCreated, gin.H{
				"provider": AdminProviderResponse{
					ID:            provider.ID,
					Name:          provider.Name,
					Slug:          provider.Slug,
					Issuer:        provider.Issuer,
					ClientID:      provider.ClientID,
					Scopes:        provider.Scopes,
					Enabled:       provider.Enabled,
					AutoProvision: provider.AutoProvision,
					CreatedAt:     provider.CreatedAt.Format(time.RFC3339),
				},
				"warning": "Provider created but failed to initialize: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, AdminProviderResponse{
		ID:            provider.ID,
		Name:          provider.Name,
		Slug:          provider.Slug,
		Issuer:        provider.Issuer,
		ClientID:      provider.ClientID,
		Scopes:        provider.Scopes,
		Enabled:       provider.Enabled,
		AutoProvision: provider.AutoProvision,
		CreatedAt:     provider.CreatedAt.Format(time.RFC3339),
	})
}

// UpdateProviderRequest represents a request to update an OIDC provider
type UpdateProviderRequest struct {
	Name          *string `json:"name"`
	Issuer        *string `json:"issuer"`
	ClientID      *string `json:"client_id"`
	ClientSecret  *string `json:"client_secret"`
	Scopes        *string `json:"scopes"`
	Enabled       *bool   `json:"enabled"`
	AutoProvision *bool   `json:"auto_provision"`
}

// UpdateProvider updates an OIDC provider
func (h *Handler) UpdateProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	var provider models.OIDCProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	var req UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Issuer != nil {
		updates["issuer"] = *req.Issuer
	}
	if req.ClientID != nil {
		updates["client_id"] = *req.ClientID
	}
	if req.ClientSecret != nil {
		updates["client_secret"] = *req.ClientSecret
	}
	if req.Scopes != nil {
		updates["scopes"] = *req.Scopes
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.AutoProvision != nil {
		updates["auto_provision"] = *req.AutoProvision
	}

	if len(updates) > 0 {
		h.db.Model(&provider).Updates(updates)
	}

	// Reload provider
	h.db.First(&provider, id)

	// Reinitialize provider
	h.mu.Lock()
	delete(h.providers, provider.ID)
	h.mu.Unlock()

	if provider.Enabled {
		h.initProvider(provider)
	}

	c.JSON(http.StatusOK, AdminProviderResponse{
		ID:            provider.ID,
		Name:          provider.Name,
		Slug:          provider.Slug,
		Issuer:        provider.Issuer,
		ClientID:      provider.ClientID,
		Scopes:        provider.Scopes,
		Enabled:       provider.Enabled,
		AutoProvision: provider.AutoProvision,
		CreatedAt:     provider.CreatedAt.Format(time.RFC3339),
	})
}

// DeleteProvider deletes an OIDC provider
func (h *Handler) DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	var provider models.OIDCProvider
	if err := h.db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	// Remove from cache
	h.mu.Lock()
	delete(h.providers, provider.ID)
	h.mu.Unlock()

	// Delete provider and associated identities
	h.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("provider_id = ?", provider.ID).Delete(&models.OIDCIdentity{})
		return tx.Delete(&provider).Error
	})

	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted"})
}

// RegisterRoutes registers public OIDC routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", h.ListProviders)
	rg.POST("/providers/:slug/auth", h.GetAuthURL)
	rg.GET("/callback", h.Callback)
}

// RegisterAdminRoutes registers admin OIDC routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", h.ListProvidersAdmin)
	rg.POST("/providers", h.CreateProvider)
	rg.PUT("/providers/:id", h.UpdateProvider)
	rg.DELETE("/providers/:id", h.DeleteProvider)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}
