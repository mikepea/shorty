package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
)

// MemberResponse represents a group member in API responses
type MemberResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// AddMemberRequest represents a request to add a member
type AddMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin member"`
}

// UpdateMemberRequest represents a request to update a member's role
type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

// ListMembers returns all members of a group
func (h *Handler) ListMembers(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check membership
	if err := h.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&models.GroupMembership{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var memberships []models.GroupMembership
	if err := h.db.Preload("User").Where("group_id = ?", groupID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	members := make([]MemberResponse, len(memberships))
	for i, m := range memberships {
		members[i] = MemberResponse{
			ID:    m.User.ID,
			Email: m.User.Email,
			Name:  m.User.Name,
			Role:  string(m.Role),
		}
	}

	c.JSON(http.StatusOK, members)
}

// AddMember adds a user to a group (admin only)
func (h *Handler) AddMember(c *gin.Context) {
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

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var targetUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already a member
	var existingMembership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ?", targetUser.ID, groupID).First(&existingMembership).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member"})
		return
	}

	// Create membership
	membership := models.GroupMembership{
		UserID:  targetUser.ID,
		GroupID: uint(groupID),
		Role:    models.GroupRole(req.Role),
	}

	if err := h.db.Create(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusCreated, MemberResponse{
		ID:    targetUser.ID,
		Email: targetUser.Email,
		Name:  targetUser.Name,
		Role:  req.Role,
	})
}

// UpdateMember updates a member's role (admin only)
func (h *Handler) UpdateMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	memberID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND group_id = ? AND role = ?", userID, groupID, models.GroupRoleAdmin).First(&models.GroupMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find membership
	var membership models.GroupMembership
	if err := h.db.Preload("User").Where("user_id = ? AND group_id = ?", memberID, groupID).First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// Update role
	membership.Role = models.GroupRole(req.Role)
	if err := h.db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update member"})
		return
	}

	c.JSON(http.StatusOK, MemberResponse{
		ID:    membership.User.ID,
		Email: membership.User.Email,
		Name:  membership.User.Name,
		Role:  string(membership.Role),
	})
}

// RemoveMember removes a user from a group (admin only)
func (h *Handler) RemoveMember(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	memberID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check admin membership
	if err := h.db.Where("user_id = ? AND group_id = ? AND role = ?", userID, groupID, models.GroupRoleAdmin).First(&models.GroupMembership{}).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Prevent removing self if only admin
	if userID == uint(memberID) {
		var adminCount int64
		h.db.Model(&models.GroupMembership{}).Where("group_id = ? AND role = ?", groupID, models.GroupRoleAdmin).Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove the last admin"})
			return
		}
	}

	// Delete membership
	result := h.db.Where("user_id = ? AND group_id = ?", memberID, groupID).Delete(&models.GroupMembership{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed"})
}

// RegisterMemberRoutes registers member management routes
func (h *Handler) RegisterMemberRoutes(rg *gin.RouterGroup) {
	rg.GET("/:id/members", h.ListMembers)
	rg.POST("/:id/members", h.AddMember)
	rg.PUT("/:id/members/:userId", h.UpdateMember)
	rg.DELETE("/:id/members/:userId", h.RemoveMember)
}
