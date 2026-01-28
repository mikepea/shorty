package organizations

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

// Handler handles organization-related requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new organizations handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateOrgRequest represents the request to create an organization
type CreateOrgRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	Slug string `json:"slug" binding:"required,min=1,max=50"`
}

// UpdateOrgRequest represents the request to update an organization
type UpdateOrgRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=100"`
}

// OrgResponse represents an organization in API responses
type OrgResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	IsGlobal    bool   `json:"is_global"`
	Role        string `json:"role,omitempty"`     // User's role in this org
	MemberCount int    `json:"member_count,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// MemberResponse represents a member in API responses
type MemberResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// AddMemberRequest represents the request to add a member
type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin member"`
}

// UpdateMemberRequest represents the request to update a member's role
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

// validateSlug checks if an organization slug is valid and available
func (h *Handler) validateSlug(slug string, excludeID uint) error {
	if slug == "" {
		return &ValidationError{"Slug is required"}
	}

	// Check format (lowercase alphanumeric with hyphens, no leading/trailing hyphens)
	if !slugRegex.MatchString(slug) {
		return &ValidationError{"Slug must contain only lowercase letters, numbers, and hyphens (no leading/trailing hyphens)"}
	}

	// Check reserved slugs
	reserved := []string{"api", "health", "admin", "login", "logout", "register", "auth", "shorty-global"}
	for _, r := range reserved {
		if strings.EqualFold(slug, r) {
			return &ValidationError{"This slug is reserved"}
		}
	}

	// Check uniqueness
	var existing models.Organization
	query := h.db.Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.First(&existing).Error; err == nil {
		return &ValidationError{"This slug is already taken"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// List returns all organizations the current user is a member of
// @Summary List organizations
// @Description Get all organizations the current user is a member of
// @Tags organizations
// @Produce json
// @Success 200 {array} OrgResponse
// @Security BearerAuth
// @Router /organizations [get]
func (h *Handler) List(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var memberships []models.OrganizationMembership
	if err := h.db.Preload("Organization").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	orgs := make([]OrgResponse, len(memberships))
	for i, m := range memberships {
		var memberCount int64
		h.db.Model(&models.OrganizationMembership{}).Where("organization_id = ?", m.OrganizationID).Count(&memberCount)

		orgs[i] = OrgResponse{
			ID:          m.Organization.ID,
			Name:        m.Organization.Name,
			Slug:        m.Organization.Slug,
			IsGlobal:    m.Organization.IsGlobal,
			Role:        string(m.Role),
			MemberCount: int(memberCount),
			CreatedAt:   m.Organization.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, orgs)
}

// Create creates a new organization and adds the creator as admin
// @Summary Create an organization
// @Description Create a new organization with the current user as admin
// @Tags organizations
// @Accept json
// @Produce json
// @Param request body CreateOrgRequest true "Organization details"
// @Success 201 {object} OrgResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Security BearerAuth
// @Router /organizations [post]
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize and validate slug
	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	if err := h.validateSlug(slug, 0); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create organization in a transaction
	var org models.Organization
	err := h.db.Transaction(func(tx *gorm.DB) error {
		org = models.Organization{
			Name: strings.TrimSpace(req.Name),
			Slug: slug,
		}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		// Add creator as admin
		membership := models.OrganizationMembership{
			OrganizationID: org.ID,
			UserID:         userID,
			Role:           models.OrgRoleAdmin,
		}
		return tx.Create(&membership).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, OrgResponse{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		IsGlobal:    org.IsGlobal,
		Role:        string(models.OrgRoleAdmin),
		MemberCount: 1,
		CreatedAt:   org.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Get returns a specific organization
// @Summary Get an organization
// @Description Get details of a specific organization
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} OrgResponse
// @Failure 404 {object} map[string]string "Organization not found"
// @Security BearerAuth
// @Router /organizations/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check membership
	var membership models.OrganizationMembership
	if err := h.db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var memberCount int64
	h.db.Model(&models.OrganizationMembership{}).Where("organization_id = ?", orgID).Count(&memberCount)

	c.JSON(http.StatusOK, OrgResponse{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		IsGlobal:    org.IsGlobal,
		Role:        string(membership.Role),
		MemberCount: int(memberCount),
		CreatedAt:   org.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Update updates an organization (admin only)
// @Summary Update an organization
// @Description Update an organization (requires admin role in org)
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param request body UpdateOrgRequest true "Updated organization details"
// @Success 200 {object} OrgResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /organizations/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check admin membership
	var membership models.OrganizationMembership
	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&membership).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	// Cannot modify global organization name
	if org.IsGlobal {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify the global organization"})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		org.Name = strings.TrimSpace(req.Name)
	}

	if err := h.db.Save(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	var memberCount int64
	h.db.Model(&models.OrganizationMembership{}).Where("organization_id = ?", orgID).Count(&memberCount)

	c.JSON(http.StatusOK, OrgResponse{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		IsGlobal:    org.IsGlobal,
		Role:        string(membership.Role),
		MemberCount: int(memberCount),
		CreatedAt:   org.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// Delete deletes an organization (admin only, soft delete)
// @Summary Delete an organization
// @Description Delete an organization (requires admin role, soft delete)
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]string "Organization deleted"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /organizations/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Check if this is the global organization
	var org models.Organization
	if err := h.db.First(&org, orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	if org.IsGlobal {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete the global organization"})
		return
	}

	// Soft delete organization (cascades preserve all data)
	if err := h.db.Delete(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted"})
}

// ListMembers returns all members of an organization
// @Summary List organization members
// @Description Get all members of an organization
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {array} MemberResponse
// @Failure 404 {object} map[string]string "Organization not found"
// @Security BearerAuth
// @Router /organizations/{id}/members [get]
func (h *Handler) ListMembers(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check membership
	if err := h.db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var memberships []models.OrganizationMembership
	if err := h.db.Preload("User").Where("organization_id = ?", orgID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	members := make([]MemberResponse, len(memberships))
	for i, m := range memberships {
		members[i] = MemberResponse{
			ID:        m.ID,
			UserID:    m.UserID,
			Email:     m.User.Email,
			Name:      m.User.Name,
			Role:      string(m.Role),
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, members)
}

// AddMember adds a user to an organization (admin only)
// @Summary Add a member to an organization
// @Description Add a user to an organization by email (requires admin role)
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param request body AddMemberRequest true "Member details"
// @Success 201 {object} MemberResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /organizations/{id}/members [post]
func (h *Handler) AddMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already a member
	var existing models.OrganizationMembership
	if err := h.db.Where("organization_id = ? AND user_id = ?", orgID, user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member"})
		return
	}

	// Create membership
	membership := models.OrganizationMembership{
		OrganizationID: uint(orgID),
		UserID:         user.ID,
		Role:           models.OrgRole(req.Role),
	}
	if err := h.db.Create(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusCreated, MemberResponse{
		ID:        membership.ID,
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(membership.Role),
		CreatedAt: membership.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateMember updates a member's role (admin only)
// @Summary Update a member's role
// @Description Update a member's role in an organization (requires admin role)
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param userId path int true "User ID"
// @Param request body UpdateMemberRequest true "Updated role"
// @Success 200 {object} MemberResponse
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /organizations/{id}/members/{userId} [put]
func (h *Handler) UpdateMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find membership
	var membership models.OrganizationMembership
	if err := h.db.Preload("User").Where("organization_id = ? AND user_id = ?", orgID, targetUserID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// Cannot demote yourself if you're the only admin
	if userID == uint(targetUserID) && req.Role == "member" {
		var adminCount int64
		h.db.Model(&models.OrganizationMembership{}).Where("organization_id = ? AND role = ?", orgID, models.OrgRoleAdmin).Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot demote the only admin"})
			return
		}
	}

	membership.Role = models.OrgRole(req.Role)
	if err := h.db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member"})
		return
	}

	c.JSON(http.StatusOK, MemberResponse{
		ID:        membership.ID,
		UserID:    membership.UserID,
		Email:     membership.User.Email,
		Name:      membership.User.Name,
		Role:      string(membership.Role),
		CreatedAt: membership.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// RemoveMember removes a member from an organization (admin only)
// @Summary Remove a member from an organization
// @Description Remove a member from an organization (requires admin role)
// @Tags organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]string "Member removed"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /organizations/{id}/members/{userId} [delete]
func (h *Handler) RemoveMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check admin membership (or self-removal)
	if userID != uint(targetUserID) {
		if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
	}

	// Find membership
	var membership models.OrganizationMembership
	if err := h.db.Where("organization_id = ? AND user_id = ?", orgID, targetUserID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// Cannot remove yourself if you're the only admin
	if userID == uint(targetUserID) && membership.Role == models.OrgRoleAdmin {
		var adminCount int64
		h.db.Model(&models.OrganizationMembership{}).Where("organization_id = ? AND role = ?", orgID, models.OrgRoleAdmin).Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove the only admin"})
			return
		}
	}

	if err := h.db.Delete(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// RegisterRoutes registers organization routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.GET("/:id", h.Get)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

// RegisterMemberRoutes registers member management routes
func (h *Handler) RegisterMemberRoutes(rg *gin.RouterGroup) {
	rg.GET("/:id/members", h.ListMembers)
	rg.POST("/:id/members", h.AddMember)
	rg.PUT("/:id/members/:userId", h.UpdateMember)
	rg.DELETE("/:id/members/:userId", h.RemoveMember)
}
