package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles admin requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new admin handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// UserResponse represents user data in admin responses
type UserResponse struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	SystemRole string `json:"system_role"`
	CreatedAt  string `json:"created_at"`
	LinkCount  int64  `json:"link_count"`
	GroupCount int64  `json:"group_count"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name       *string `json:"name"`
	SystemRole *string `json:"system_role"`
}

// StatsResponse represents system statistics
type StatsResponse struct {
	TotalUsers      int64 `json:"total_users"`
	TotalLinks      int64 `json:"total_links"`
	TotalGroups     int64 `json:"total_groups"`
	TotalTags       int64 `json:"total_tags"`
	TotalClicks     int64 `json:"total_clicks"`
	PublicLinks     int64 `json:"public_links"`
	PrivateLinks    int64 `json:"private_links"`
	UnreadLinks     int64 `json:"unread_links"`
	AdminUsers      int64 `json:"admin_users"`
	ActiveAPIKeys   int64 `json:"active_api_keys"`
}

// ListUsers returns all users (admin only)
func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User

	query := h.db.Order("created_at DESC")

	// Optional search by email or name
	if search := c.Query("q"); search != "" {
		query = query.Where("email LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Optional filter by role
	if role := c.Query("role"); role != "" {
		query = query.Where("system_role = ?", role)
	}

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		var linkCount, groupCount int64
		h.db.Model(&models.Link{}).Where("created_by_id = ?", user.ID).Count(&linkCount)
		h.db.Model(&models.GroupMembership{}).Where("user_id = ?", user.ID).Count(&groupCount)

		responses[i] = UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Name:       user.Name,
			SystemRole: string(user.SystemRole),
			CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LinkCount:  linkCount,
			GroupCount: groupCount,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// GetUser returns a single user by ID (admin only)
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var linkCount, groupCount int64
	h.db.Model(&models.Link{}).Where("created_by_id = ?", user.ID).Count(&linkCount)
	h.db.Model(&models.GroupMembership{}).Where("user_id = ?", user.ID).Count(&groupCount)

	c.JSON(http.StatusOK, UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		SystemRole: string(user.SystemRole),
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LinkCount:  linkCount,
		GroupCount: groupCount,
	})
}

// UpdateUser updates a user's profile (admin only)
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent admin from demoting themselves
	currentUserID, _ := auth.GetUserID(c)
	if uint(id) == currentUserID && req.SystemRole != nil && *req.SystemRole != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot demote yourself"})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.SystemRole != nil {
		if *req.SystemRole != "admin" && *req.SystemRole != "user" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system role"})
			return
		}
		updates["system_role"] = *req.SystemRole
	}

	if len(updates) > 0 {
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Reload user
	h.db.First(&user, id)

	var linkCount, groupCount int64
	h.db.Model(&models.Link{}).Where("created_by_id = ?", user.ID).Count(&linkCount)
	h.db.Model(&models.GroupMembership{}).Where("user_id = ?", user.ID).Count(&groupCount)

	c.JSON(http.StatusOK, UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		SystemRole: string(user.SystemRole),
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LinkCount:  linkCount,
		GroupCount: groupCount,
	})
}

// DeleteUser soft-deletes a user (admin only)
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent admin from deleting themselves
	currentUserID, _ := auth.GetUserID(c)
	if uint(id) == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user and related data in a transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Delete API keys
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.APIKey{}).Error; err != nil {
			return err
		}
		// Delete group memberships
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.GroupMembership{}).Error; err != nil {
			return err
		}
		// Delete links
		if err := tx.Where("created_by_id = ?", user.ID).Delete(&models.Link{}).Error; err != nil {
			return err
		}
		// Delete user
		return tx.Delete(&user).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetStats returns system-wide statistics (admin only)
func (h *Handler) GetStats(c *gin.Context) {
	var stats StatsResponse

	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
	h.db.Model(&models.Link{}).Count(&stats.TotalLinks)
	h.db.Model(&models.Group{}).Count(&stats.TotalGroups)
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)
	h.db.Model(&models.APIKey{}).Count(&stats.ActiveAPIKeys)

	h.db.Model(&models.Link{}).Where("is_public = ?", true).Count(&stats.PublicLinks)
	h.db.Model(&models.Link{}).Where("is_public = ?", false).Count(&stats.PrivateLinks)
	h.db.Model(&models.Link{}).Where("is_unread = ?", true).Count(&stats.UnreadLinks)
	h.db.Model(&models.User{}).Where("system_role = ?", "admin").Count(&stats.AdminUsers)

	// Sum of all click counts
	h.db.Model(&models.Link{}).Select("COALESCE(SUM(click_count), 0)").Scan(&stats.TotalClicks)

	c.JSON(http.StatusOK, stats)
}

// RegisterRoutes registers admin routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/stats", h.GetStats)
	rg.GET("/users", h.ListUsers)
	rg.GET("/users/:id", h.GetUser)
	rg.PUT("/users/:id", h.UpdateUser)
	rg.DELETE("/users/:id", h.DeleteUser)
}
