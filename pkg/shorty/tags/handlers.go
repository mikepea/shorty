package tags

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles tag-related requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new tags handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LinkCount int    `json:"link_count,omitempty"`
}

// SetTagsRequest represents the request to set tags on a link
type SetTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

// getUserGroupIDs returns all group IDs the user is a member of
func (h *Handler) getUserGroupIDs(userID string) ([]string, error) {
	var memberships []models.GroupMembership
	if err := h.db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return nil, err
	}

	groupIDs := make([]string, len(memberships))
	for i, m := range memberships {
		groupIDs[i] = m.GroupID
	}
	return groupIDs, nil
}

// checkGroupMembership verifies the user is a member of the group
func (h *Handler) checkGroupMembership(userID string, groupID string) error {
	var membership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&membership).Error; err != nil {
		return err
	}
	return nil
}

// List returns all tags used across the user's groups
// @Summary List tags
// @Description Get all tags used across the user's groups
// @Tags tags
// @Produce json
// @Success 200 {array} TagResponse
// @Security BearerAuth
// @Router /tags [get]
func (h *Handler) List(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	groupIDs, err := h.getUserGroupIDs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	if len(groupIDs) == 0 {
		c.JSON(http.StatusOK, []TagResponse{})
		return
	}

	// Get tags with link counts for user's groups
	type tagWithCount struct {
		ID        string
		Name      string
		LinkCount int
	}

	var results []tagWithCount
	err = h.db.Table("tags").
		Select("tags.id, tags.name, COUNT(DISTINCT links.id) as link_count").
		Joins("INNER JOIN link_tags ON tags.id = link_tags.tag_id").
		Joins("INNER JOIN links ON link_tags.link_id = links.id AND links.group_id IN ? AND links.deleted_at IS NULL", groupIDs).
		Where("tags.deleted_at IS NULL").
		Group("tags.id").
		Order("link_count DESC").
		Find(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	tags := make([]TagResponse, len(results))
	for i, r := range results {
		tags[i] = TagResponse{
			ID:        r.ID,
			Name:      r.Name,
			LinkCount: r.LinkCount,
		}
	}

	c.JSON(http.StatusOK, tags)
}

// ListByGroup returns all tags used in a specific group
// @Summary List tags in a group
// @Description Get all tags used in a specific group
// @Tags tags
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {array} TagResponse
// @Failure 404 {object} map[string]string "Group not found"
// @Security BearerAuth
// @Router /groups/{id}/tags [get]
func (h *Handler) ListByGroup(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID := c.Param("id")

	// Check membership
	if err := h.checkGroupMembership(userID, groupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	type tagWithCount struct {
		ID        string
		Name      string
		LinkCount int
	}

	var results []tagWithCount
	err := h.db.Table("tags").
		Select("tags.id, tags.name, COUNT(DISTINCT links.id) as link_count").
		Joins("INNER JOIN link_tags ON tags.id = link_tags.tag_id").
		Joins("INNER JOIN links ON link_tags.link_id = links.id AND links.group_id = ? AND links.deleted_at IS NULL", groupID).
		Where("tags.deleted_at IS NULL").
		Group("tags.id").
		Order("link_count DESC").
		Find(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	tags := make([]TagResponse, len(results))
	for i, r := range results {
		tags[i] = TagResponse{
			ID:        r.ID,
			Name:      r.Name,
			LinkCount: r.LinkCount,
		}
	}

	c.JSON(http.StatusOK, tags)
}

// GetLinkTags returns tags for a specific link
// @Summary Get link tags
// @Description Get all tags for a specific link
// @Tags tags
// @Produce json
// @Param slug path string true "Link slug"
// @Success 200 {array} TagResponse
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug}/tags [get]
func (h *Handler) GetLinkTags(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	slug := c.Param("slug")

	var link models.Link
	if err := h.db.Preload("Tags").Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Check access
	if !link.IsPublic {
		if err := h.checkGroupMembership(userID, link.GroupID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
	}

	tags := make([]TagResponse, len(link.Tags))
	for i, t := range link.Tags {
		tags[i] = TagResponse{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	c.JSON(http.StatusOK, tags)
}

// SetLinkTags sets the tags for a link (replaces existing tags)
// @Summary Set link tags
// @Description Replace all tags on a link
// @Tags tags
// @Accept json
// @Produce json
// @Param slug path string true "Link slug"
// @Param request body SetTagsRequest true "Tags to set"
// @Success 200 {array} TagResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug}/tags [put]
func (h *Handler) SetLinkTags(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	slug := c.Param("slug")

	var link models.Link
	if err := h.db.Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Check membership
	if err := h.checkGroupMembership(userID, link.GroupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	var req SetTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create tags
	var tags []models.Tag
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}

		var tag models.Tag
		// Try to find existing tag
		if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
			// Create new tag
			tag = models.Tag{Name: tagName}
			if err := h.db.Create(&tag).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
				return
			}
		}
		tags = append(tags, tag)
	}

	// Replace link's tags
	if err := h.db.Model(&link).Association("Tags").Replace(tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags"})
		return
	}

	// Return updated tags
	tagResponses := make([]TagResponse, len(tags))
	for i, t := range tags {
		tagResponses[i] = TagResponse{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	c.JSON(http.StatusOK, tagResponses)
}

// AddLinkTag adds a single tag to a link
// @Summary Add tag to link
// @Description Add a single tag to a link
// @Tags tags
// @Produce json
// @Param slug path string true "Link slug"
// @Param tag path string true "Tag name"
// @Success 200 {object} TagResponse
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug}/tags/{tag} [post]
func (h *Handler) AddLinkTag(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	slug := c.Param("slug")
	tagName := c.Param("tag")

	var link models.Link
	if err := h.db.Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Check membership
	if err := h.checkGroupMembership(userID, link.GroupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Get or create tag
	var tag models.Tag
	if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		tag = models.Tag{Name: tagName}
		if err := h.db.Create(&tag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
			return
		}
	}

	// Add tag to link
	if err := h.db.Model(&link).Association("Tags").Append(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag"})
		return
	}

	c.JSON(http.StatusOK, TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	})
}

// RemoveLinkTag removes a tag from a link
// @Summary Remove tag from link
// @Description Remove a tag from a link
// @Tags tags
// @Produce json
// @Param slug path string true "Link slug"
// @Param tag path string true "Tag name"
// @Success 200 {object} map[string]string "Tag removed"
// @Failure 404 {object} map[string]string "Link or tag not found"
// @Security BearerAuth
// @Router /links/{slug}/tags/{tag} [delete]
func (h *Handler) RemoveLinkTag(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	slug := c.Param("slug")
	tagName := c.Param("tag")

	var link models.Link
	if err := h.db.Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Check membership
	if err := h.checkGroupMembership(userID, link.GroupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Find tag
	var tag models.Tag
	if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Remove tag from link
	if err := h.db.Model(&link).Association("Tags").Delete(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed"})
}

// RegisterRoutes registers tag routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// List all tags across user's groups
	rg.GET("/tags", h.List)

	// List tags in a specific group
	rg.GET("/groups/:id/tags", h.ListByGroup)

	// Link tag operations
	rg.GET("/links/:slug/tags", h.GetLinkTags)
	rg.PUT("/links/:slug/tags", h.SetLinkTags)
	rg.POST("/links/:slug/tags/:tag", h.AddLinkTag)
	rg.DELETE("/links/:slug/tags/:tag", h.RemoveLinkTag)
}
