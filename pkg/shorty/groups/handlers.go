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
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
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
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create group in a transaction
	var group models.Group
	err := h.db.Transaction(func(tx *gorm.DB) error {
		group = models.Group{
			Name:        req.Name,
			Description: req.Description,
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
