package links

import (
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

var slugRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Handler handles link-related requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new links handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// CreateLinkRequest represents the request to create a link
type CreateLinkRequest struct {
	URL         string `json:"url" binding:"required,url"`
	Slug        string `json:"slug" binding:"omitempty,min=1,max=50"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsUnread    bool   `json:"is_unread"`
}

// UpdateLinkRequest represents the request to update a link
type UpdateLinkRequest struct {
	URL         string `json:"url" binding:"omitempty,url"`
	Slug        string `json:"slug" binding:"omitempty,min=1,max=50"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"`
	IsUnread    *bool  `json:"is_unread"`
}

// LinkResponse represents a link in API responses
type LinkResponse struct {
	ID          uint   `json:"id"`
	GroupID     uint   `json:"group_id"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsUnread    bool   `json:"is_unread"`
	ClickCount  uint   `json:"click_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func linkToResponse(link models.Link) LinkResponse {
	return LinkResponse{
		ID:          link.ID,
		GroupID:     link.GroupID,
		Slug:        link.Slug,
		URL:         link.URL,
		Title:       link.Title,
		Description: link.Description,
		IsPublic:    link.IsPublic,
		IsUnread:    link.IsUnread,
		ClickCount:  link.ClickCount,
		CreatedAt:   link.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   link.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// validateSlug checks if a slug is valid and available (deprecated - use validateSlugForOrg)
func (h *Handler) validateSlug(slug string, excludeID uint) error {
	return h.validateSlugForOrg(slug, excludeID, 0)
}

// validateSlugForOrg checks if a slug is valid and available within an organization
func (h *Handler) validateSlugForOrg(slug string, excludeID uint, orgID uint) error {
	if slug == "" {
		return nil
	}

	// Check format
	if !slugRegex.MatchString(slug) {
		return &ValidationError{"Slug must contain only letters, numbers, hyphens, and underscores"}
	}

	// Check reserved slugs
	reserved := []string{"api", "health", "admin", "login", "logout", "register", "auth"}
	for _, r := range reserved {
		if strings.EqualFold(slug, r) {
			return &ValidationError{"This slug is reserved"}
		}
	}

	// Check uniqueness within organization
	var existing models.Link
	query := h.db.Where("organization_id = ? AND slug = ?", orgID, slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.First(&existing).Error; err == nil {
		return &ValidationError{"This slug is already taken"}
	}

	return nil
}

// generateRandomString creates a random string of given length
func generateRandomString(length int, charset string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// generateSlug creates a unique slug (deprecated - use generateSlugForOrg)
func (h *Handler) generateSlug() string {
	return h.generateSlugForOrg(0)
}

// generateSlugForOrg creates a unique slug within an organization
func (h *Handler) generateSlugForOrg(orgID uint) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	for attempts := 0; attempts < 10; attempts++ {
		slug := generateRandomString(length, charset)
		var existing models.Link
		if err := h.db.Where("organization_id = ? AND slug = ?", orgID, slug).First(&existing).Error; err != nil {
			return slug
		}
	}

	// Fallback to longer slug if short ones are exhausted
	return generateRandomString(12, charset)
}

// checkGroupMembership verifies the user is a member of the group
func (h *Handler) checkGroupMembership(userID, groupID uint) error {
	var membership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&membership).Error; err != nil {
		return err
	}
	return nil
}

// ListByGroup returns all links in a group
// @Summary List links in a group
// @Description Get all links belonging to a specific group
// @Tags links
// @Produce json
// @Param id path int true "Group ID"
// @Param is_unread query bool false "Filter by unread status"
// @Param is_public query bool false "Filter by public status"
// @Success 200 {array} LinkResponse
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 404 {object} map[string]string "Group not found"
// @Security BearerAuth
// @Router /groups/{id}/links [get]
func (h *Handler) ListByGroup(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check membership
	if err := h.checkGroupMembership(userID, uint(groupID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var links []models.Link
	query := h.db.Where("group_id = ?", groupID).Order("created_at DESC")

	// Optional filters
	if isUnread := c.Query("is_unread"); isUnread != "" {
		query = query.Where("is_unread = ?", isUnread == "true")
	}
	if isPublic := c.Query("is_public"); isPublic != "" {
		query = query.Where("is_public = ?", isPublic == "true")
	}

	if err := query.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	responses := make([]LinkResponse, len(links))
	for i, link := range links {
		responses[i] = linkToResponse(link)
	}

	c.JSON(http.StatusOK, responses)
}

// Create creates a new link in a group
// @Summary Create a link
// @Description Create a new shortened link in a group
// @Tags links
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body CreateLinkRequest true "Link details"
// @Success 201 {object} LinkResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 404 {object} map[string]string "Group not found"
// @Security BearerAuth
// @Router /groups/{id}/links [post]
func (h *Handler) Create(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Check membership
	if err := h.checkGroupMembership(userID, uint(groupID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Get the group to find its organization
	var group models.Group
	if err := h.db.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle slug - now scoped to organization
	slug := req.Slug
	if slug == "" {
		slug = h.generateSlugForOrg(group.OrganizationID)
	} else {
		if err := h.validateSlugForOrg(slug, 0, group.OrganizationID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	link := models.Link{
		OrganizationID: group.OrganizationID,
		GroupID:        uint(groupID),
		CreatedByID:    userID,
		Slug:           slug,
		URL:            req.URL,
		Title:          req.Title,
		Description:    req.Description,
		IsPublic:       req.IsPublic,
		IsUnread:       req.IsUnread,
	}

	if err := h.db.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}

	c.JSON(http.StatusCreated, linkToResponse(link))
}

// GetBySlug returns a link by its slug
// @Summary Get a link by slug
// @Description Get link details by its short slug
// @Tags links
// @Produce json
// @Param slug path string true "Link slug"
// @Success 200 {object} LinkResponse
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug} [get]
func (h *Handler) GetBySlug(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	slug := c.Param("slug")

	var link models.Link
	if err := h.db.Preload("Tags").Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Check if user has access (public or member of group)
	if !link.IsPublic {
		if err := h.checkGroupMembership(userID, link.GroupID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
	}

	c.JSON(http.StatusOK, linkToResponse(link))
}

// Update updates a link
// @Summary Update a link
// @Description Update an existing link by slug
// @Tags links
// @Accept json
// @Produce json
// @Param slug path string true "Link slug"
// @Param request body UpdateLinkRequest true "Updated link details"
// @Success 200 {object} LinkResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug} [put]
func (h *Handler) Update(c *gin.Context) {
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

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate new slug if provided
	if req.Slug != "" && req.Slug != link.Slug {
		if err := h.validateSlug(req.Slug, link.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link.Slug = req.Slug
	}

	// Update fields
	if req.URL != "" {
		link.URL = req.URL
	}
	if req.Title != "" {
		link.Title = req.Title
	}
	if req.Description != "" {
		link.Description = req.Description
	}
	if req.IsPublic != nil {
		link.IsPublic = *req.IsPublic
	}
	if req.IsUnread != nil {
		link.IsUnread = *req.IsUnread
	}

	if err := h.db.Save(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update link"})
		return
	}

	c.JSON(http.StatusOK, linkToResponse(link))
}

// Delete deletes a link
// @Summary Delete a link
// @Description Delete a link by slug
// @Tags links
// @Produce json
// @Param slug path string true "Link slug"
// @Success 200 {object} map[string]string "Link deleted"
// @Failure 404 {object} map[string]string "Link not found"
// @Security BearerAuth
// @Router /links/{slug} [delete]
func (h *Handler) Delete(c *gin.Context) {
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

	if err := h.db.Delete(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted"})
}

// Search searches links across all user's groups
// @Summary Search links
// @Description Search links across all groups the user has access to
// @Tags links
// @Produce json
// @Param q query string false "Search query (searches title, description, URL)"
// @Param is_unread query bool false "Filter by unread status"
// @Param is_public query bool false "Filter by public status"
// @Param group_id query int false "Filter by group ID"
// @Param tag query string false "Filter by tag name"
// @Param limit query int false "Max results (default 50, max 100)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} LinkResponse
// @Security BearerAuth
// @Router /links [get]
func (h *Handler) Search(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	// Get user's group IDs
	var memberships []models.GroupMembership
	if err := h.db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	groupIDs := make([]uint, len(memberships))
	for i, m := range memberships {
		groupIDs[i] = m.GroupID
	}

	if len(groupIDs) == 0 {
		c.JSON(http.StatusOK, []LinkResponse{})
		return
	}

	query := h.db.Where("group_id IN ?", groupIDs).Order("created_at DESC")

	// Search term
	if q := c.Query("q"); q != "" {
		searchTerm := "%" + q + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR url LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Filters
	if isUnread := c.Query("is_unread"); isUnread != "" {
		query = query.Where("is_unread = ?", isUnread == "true")
	}
	if isPublic := c.Query("is_public"); isPublic != "" {
		query = query.Where("is_public = ?", isPublic == "true")
	}
	if groupID := c.Query("group_id"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if tag := c.Query("tag"); tag != "" {
		query = query.Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// Pagination
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	query = query.Limit(limit)

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	query = query.Offset(offset)

	var links []models.Link
	if err := query.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search links"})
		return
	}

	responses := make([]LinkResponse, len(links))
	for i, link := range links {
		responses[i] = linkToResponse(link)
	}

	c.JSON(http.StatusOK, responses)
}

// RegisterRoutes registers link routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Group-scoped routes
	rg.GET("/groups/:id/links", h.ListByGroup)
	rg.POST("/groups/:id/links", h.Create)

	// Slug-based routes
	rg.GET("/links/:slug", h.GetBySlug)
	rg.PUT("/links/:slug", h.Update)
	rg.DELETE("/links/:slug", h.Delete)

	// Search across all groups
	rg.GET("/links", h.Search)
}
