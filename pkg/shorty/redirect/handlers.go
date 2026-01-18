package redirect

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles redirect requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new redirect handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Redirect handles short URL redirects
// Public links redirect without authentication
// Private links also redirect (the URL itself is not secret, just the metadata)
// Click count is incremented for all redirects
func (h *Handler) Redirect(c *gin.Context) {
	slug := c.Param("slug")

	// Find the link
	var link models.Link
	if err := h.db.Where("slug = ?", slug).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Increment click count (fire and forget - don't block redirect on DB update)
	go func() {
		h.db.Model(&link).Update("click_count", gorm.Expr("click_count + 1"))
	}()

	// Redirect to the target URL
	c.Redirect(http.StatusFound, link.URL)
}

// RegisterRoutes registers redirect routes on the root router
// This should be called AFTER all other routes to avoid conflicts
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Match any path that could be a slug
	// This is registered last to avoid conflicts with /api, /health, etc.
	r.GET("/:slug", h.Redirect)
}
