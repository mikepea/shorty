package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles group-related requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new groups handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateGroupRequest represents the request to create a group
type CreateGroupRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OrganizationID uint   `json:"organization_id"` // Optional - defaults to org from context or global
}

// UpdateGroupRequest represents the request to update a group
type UpdateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GroupResponse represents a group in API responses
type GroupResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role,omitempty"` // User's role in this group
	MemberCount int    `json:"member_count,omitempty"`
}

// List returns all groups the current user is a member of
// @Summary List groups
// @Description Get all groups the current user is a member of
// @Tags groups
// @Produce json
// @Success 200 {array} GroupResponse
// @Security BearerAuth
// @Router /groups [get]
func (h *Handler) List(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var memberships []models.GroupMembership
	if err := h.db.Preload("Group").Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	groups := make([]GroupResponse, len(memberships))
	for i, m := range memberships {
		var memberCount int64
		h.db.Model(&models.GroupMembership{}).Where("group_id = ?", m.GroupID).Count(&memberCount)

		groups[i] = GroupResponse{
			ID:          m.Group.ID,
			Name:        m.Group.Name,
			Description: m.Group.Description,
			Role:        string(m.Role),
			MemberCount: int(memberCount),
		}
	}

	c.JSON(http.StatusOK, groups)
}

// Create creates a new group and adds the creator as admin
// @Summary Create a group
// @Description Create a new group with the current user as admin
// @Tags groups
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Group details"
// @Success 201 {object} GroupResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Security BearerAuth
// @Router /groups [post]
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine organization ID:
	// 1. Use request body if provided
	// 2. Use context if OrgMiddleware set it
	// 3. Default to global organization
	orgID := req.OrganizationID
	if orgID == 0 {
		if ctxOrgID, ok := auth.GetOrgID(c); ok {
			orgID = ctxOrgID
		} else {
			// Fall back to global organization
			var globalOrg models.Organization
			if err := h.db.Where("is_global = ?", true).First(&globalOrg).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found"})
				return
			}
			orgID = globalOrg.ID
		}
	}

	// Verify user is a member of the target organization
	var orgMembership models.OrganizationMembership
	if err := h.db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&orgMembership).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this organization"})
		return
	}

	// Create group in a transaction
	var group models.Group
	err := h.db.Transaction(func(tx *gorm.DB) error {
		group = models.Group{
			OrganizationID: orgID,
			Name:           req.Name,
			Description:    req.Description,
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		// Add creator as admin
		membership := models.GroupMembership{
			UserID:  userID,
			GroupID: group.ID,
			Role:    models.GroupRoleAdmin,
		}
		return tx.Create(&membership).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Role:        string(models.GroupRoleAdmin),
		MemberCount: 1,
	})
}

// Get returns a specific group
// @Summary Get a group
// @Description Get details of a specific group
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} GroupResponse
// @Failure 404 {object} map[string]string "Group not found"
// @Security BearerAuth
// @Router /groups/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check membership
	var membership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var group models.Group
	if err := h.db.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var memberCount int64
	h.db.Model(&models.GroupMembership{}).Where("group_id = ?", groupID).Count(&memberCount)

	c.JSON(http.StatusOK, GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Role:        string(membership.Role),
		MemberCount: int(memberCount),
	})
}

// Update updates a group (admin only)
// @Summary Update a group
// @Description Update a group (requires admin role in group)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body UpdateGroupRequest true "Updated group details"
// @Success 200 {object} GroupResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /groups/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check admin membership
	var membership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ? AND role = ?", userID, groupID, models.GroupRoleAdmin).First(&membership).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var group models.Group
	if err := h.db.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}

	if err := h.db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group"})
		return
	}

	var memberCount int64
	h.db.Model(&models.GroupMembership{}).Where("group_id = ?", groupID).Count(&memberCount)

	c.JSON(http.StatusOK, GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Role:        string(membership.Role),
		MemberCount: int(memberCount),
	})
}

// Delete deletes a group (admin only)
// @Summary Delete a group
// @Description Delete a group (requires admin role in group)
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]string "Group deleted"
// @Failure 403 {object} map[string]string "Admin access required"
// @Security BearerAuth
// @Router /groups/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND group_id = ? AND role = ?", userID, groupID, models.GroupRoleAdmin).First(&models.GroupMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Delete group (cascades to memberships via soft delete)
	if err := h.db.Delete(&models.Group{}, groupID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

// RegisterRoutes registers group routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.POST("", h.Create)
	rg.GET("/:id", h.Get)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}
