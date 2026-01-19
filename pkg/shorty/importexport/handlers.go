package importexport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/auth"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles import/export requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new import/export handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// PinboardBookmark represents a bookmark in Pinboard JSON format
type PinboardBookmark struct {
	Href        string `json:"href"`
	Description string `json:"description"`
	Extended    string `json:"extended"`
	Tags        string `json:"tags"`
	Time        string `json:"time"`
	Shared      string `json:"shared"`
	ToRead      string `json:"toread"`
	Meta        string `json:"meta,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

// ImportRequest represents an import request
type ImportRequest struct {
	GroupID   uint               `json:"group_id" binding:"required"`
	Bookmarks []PinboardBookmark `json:"bookmarks" binding:"required"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}

// ExportBookmark represents a bookmark for export
type ExportBookmark struct {
	Href        string `json:"href"`
	Description string `json:"description"`
	Extended    string `json:"extended"`
	Tags        string `json:"tags"`
	Time        string `json:"time"`
	Shared      string `json:"shared"`
	ToRead      string `json:"toread"`
}

// checkGroupMembership verifies the user is a member of the group
func (h *Handler) checkGroupMembership(userID, groupID uint) error {
	var membership models.GroupMembership
	if err := h.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&membership).Error; err != nil {
		return err
	}
	return nil
}

// getUserGroupIDs returns all group IDs the user is a member of
func (h *Handler) getUserGroupIDs(userID uint) ([]uint, error) {
	var memberships []models.GroupMembership
	if err := h.db.Where("user_id = ?", userID).Find(&memberships).Error; err != nil {
		return nil, err
	}

	groupIDs := make([]uint, len(memberships))
	for i, m := range memberships {
		groupIDs[i] = m.GroupID
	}
	return groupIDs, nil
}

// generateSlug generates a unique slug for a link
func (h *Handler) generateSlug() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	for attempts := 0; attempts < 10; attempts++ {
		slug := make([]byte, length)
		for i := range slug {
			slug[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			time.Sleep(time.Nanosecond)
		}

		// Check uniqueness
		var count int64
		h.db.Model(&models.Link{}).Where("slug = ?", string(slug)).Count(&count)
		if count == 0 {
			return string(slug), nil
		}
	}

	return "", nil
}

// Import imports bookmarks from Pinboard JSON format
func (h *Handler) Import(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check group membership
	if err := h.checkGroupMembership(userID, req.GroupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	result := ImportResult{
		Errors: []string{},
	}

	for i, bookmark := range req.Bookmarks {
		// Parse time
		var createdAt time.Time
		if bookmark.Time != "" {
			parsed, err := time.Parse(time.RFC3339, bookmark.Time)
			if err != nil {
				// Try alternative format
				parsed, err = time.Parse("2006-01-02T15:04:05Z", bookmark.Time)
				if err != nil {
					result.Errors = append(result.Errors, "bookmark "+strconv.Itoa(i)+": invalid time format")
					result.Skipped++
					continue
				}
			}
			createdAt = parsed
		} else {
			createdAt = time.Now()
		}

		// Generate slug
		slug, err := h.generateSlug()
		if err != nil || slug == "" {
			result.Errors = append(result.Errors, "bookmark "+strconv.Itoa(i)+": failed to generate slug")
			result.Skipped++
			continue
		}

		// Determine visibility
		isPublic := bookmark.Shared == "yes"

		// Determine unread status
		isUnread := bookmark.ToRead == "yes"

		// Create link
		link := models.Link{
			GroupID:     req.GroupID,
			CreatedByID: userID,
			Slug:        slug,
			URL:         bookmark.Href,
			Title:       bookmark.Description,
			Description: bookmark.Extended,
			IsPublic:    isPublic,
			IsUnread:    isUnread,
		}
		link.CreatedAt = createdAt

		if err := h.db.Create(&link).Error; err != nil {
			result.Errors = append(result.Errors, "bookmark "+strconv.Itoa(i)+": "+err.Error())
			result.Skipped++
			continue
		}

		// Explicitly update boolean fields to override GORM defaults
		// GORM applies defaults when values are zero values (false for bools)
		h.db.Model(&link).Updates(map[string]interface{}{
			"is_public": isPublic,
			"is_unread": isUnread,
		})

		// Handle tags
		if bookmark.Tags != "" {
			tagNames := strings.Fields(bookmark.Tags)
			var tags []models.Tag

			for _, tagName := range tagNames {
				tagName = strings.TrimSpace(tagName)
				if tagName == "" {
					continue
				}

				var tag models.Tag
				if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
					// Create new tag
					tag = models.Tag{Name: tagName}
					if err := h.db.Create(&tag).Error; err != nil {
						continue
					}
				}
				tags = append(tags, tag)
			}

			if len(tags) > 0 {
				h.db.Model(&link).Association("Tags").Append(tags)
			}
		}

		result.Imported++
	}

	c.JSON(http.StatusOK, result)
}

// Export exports bookmarks to Pinboard JSON format
func (h *Handler) Export(c *gin.Context) {
	userID, _ := auth.GetUserID(c)

	// Get optional group_id parameter
	groupIDStr := c.Query("group_id")
	var groupIDs []uint

	if groupIDStr != "" {
		groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
			return
		}

		// Check membership
		if err := h.checkGroupMembership(userID, uint(groupID)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}

		groupIDs = []uint{uint(groupID)}
	} else {
		// Export from all user's groups
		var err error
		groupIDs, err = h.getUserGroupIDs(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
			return
		}
	}

	if len(groupIDs) == 0 {
		c.JSON(http.StatusOK, []ExportBookmark{})
		return
	}

	// Fetch links with tags
	var links []models.Link
	if err := h.db.Preload("Tags").Where("group_id IN ?", groupIDs).Order("created_at DESC").Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	// Convert to Pinboard format
	bookmarks := make([]ExportBookmark, len(links))
	for i, link := range links {
		// Convert tags to space-separated string
		tagNames := make([]string, len(link.Tags))
		for j, tag := range link.Tags {
			tagNames[j] = tag.Name
		}

		shared := "no"
		if link.IsPublic {
			shared = "yes"
		}

		toread := "no"
		if link.IsUnread {
			toread = "yes"
		}

		bookmarks[i] = ExportBookmark{
			Href:        link.URL,
			Description: link.Title,
			Extended:    link.Description,
			Tags:        strings.Join(tagNames, " "),
			Time:        link.CreatedAt.Format(time.RFC3339),
			Shared:      shared,
			ToRead:      toread,
		}
	}

	// Set content disposition for download
	if c.Query("download") == "true" {
		c.Header("Content-Disposition", "attachment; filename=shorty-export.json")
	}

	c.JSON(http.StatusOK, bookmarks)
}

// ExportSingle exports a single link to Pinboard JSON format
func (h *Handler) ExportSingle(c *gin.Context) {
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

	// Convert tags to space-separated string
	tagNames := make([]string, len(link.Tags))
	for i, tag := range link.Tags {
		tagNames[i] = tag.Name
	}

	shared := "no"
	if link.IsPublic {
		shared = "yes"
	}

	toread := "no"
	if link.IsUnread {
		toread = "yes"
	}

	bookmark := ExportBookmark{
		Href:        link.URL,
		Description: link.Title,
		Extended:    link.Description,
		Tags:        strings.Join(tagNames, " "),
		Time:        link.CreatedAt.Format(time.RFC3339),
		Shared:      shared,
		ToRead:      toread,
	}

	c.JSON(http.StatusOK, bookmark)
}

// RegisterRoutes registers import/export routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/import", h.Import)
	rg.GET("/export", h.Export)
	rg.GET("/export/:slug", h.ExportSingle)
}

// Ensure json is used (for compile check)
var _ = json.Marshal
