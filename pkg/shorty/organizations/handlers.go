package organizations

import (
	"net/http"
	"regexp"
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	IsGlobal    bool   `json:"is_global"`
	Role        string `json:"role,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// MemberResponse represents a member in API responses
type MemberResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
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
func (h *Handler) validateSlug(slug string, excludeID string) error {
	if slug == "" {
		return &ValidationError{"Slug is required"}
	}

	if !slugRegex.MatchString(slug) {
		return &ValidationError{"Slug must contain only lowercase letters, numbers, and hyphens (no leading/trailing hyphens)"}
	}

	reserved := []string{"api", "health", "admin", "login", "logout", "register", "auth", "shorty-global"}
	for _, r := range reserved {
		if strings.EqualFold(slug, r) {
			return &ValidationError{"This slug is reserved"}
		}
	}

	var existing models.Organization
	query := h.db.Where("slug = ?", slug)
	if excludeID != "" {
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
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := strings.ToLower(strings.TrimSpace(req.Slug))
	if err := h.validateSlug(slug, ""); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var org models.Organization
	err := h.db.Transaction(func(tx *gorm.DB) error {
		org = models.Organization{
			Name: strings.TrimSpace(req.Name),
			Slug: slug,
		}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

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
func (h *Handler) Get(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var membership models.OrganizationMembership
	if err := h.db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, "id = ?", orgID).Error; err != nil {
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
func (h *Handler) Update(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

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
	if err := h.db.First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	if org.IsGlobal {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify the global organization"})
		return
	}

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
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var org models.Organization
	if err := h.db.First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	if org.IsGlobal {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete the global organization"})
		return
	}

	if err := h.db.Delete(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted"})
}

// ListMembers returns all members of an organization
func (h *Handler) ListMembers(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

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
func (h *Handler) AddMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existing models.OrganizationMembership
	if err := h.db.Where("organization_id = ? AND user_id = ?", orgID, user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member"})
		return
	}

	membership := models.OrganizationMembership{
		OrganizationID: orgID,
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
func (h *Handler) UpdateMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var membership models.OrganizationMembership
	if err := h.db.Preload("User").Where("organization_id = ? AND user_id = ?", orgID, targetUserID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if userID == targetUserID && req.Role == "member" {
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
func (h *Handler) RemoveMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	targetUserID := c.Param("userId")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if userID != targetUserID {
		if err := h.db.Where("user_id = ? AND organization_id = ? AND role = ?", userID, orgID, models.OrgRoleAdmin).First(&models.OrganizationMembership{}).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
	}

	var membership models.OrganizationMembership
	if err := h.db.Where("organization_id = ? AND user_id = ?", orgID, targetUserID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	if userID == targetUserID && membership.Role == models.OrgRoleAdmin {
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
