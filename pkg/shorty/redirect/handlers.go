package redirect

import (
	"net/http"
	"strings"

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

// resolveOrgFromHost looks up an organization by the request's Host header.
// If no matching domain is found, returns the global organization ID.
// Returns 0 if neither can be found (shouldn't happen if DB is properly seeded).
func (h *Handler) resolveOrgFromHost(c *gin.Context) uint {
	host := c.Request.Host

	// Remove port if present (e.g., "localhost:8080" -> "localhost")
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	// Look up domain in OrganizationDomain table
	var domain models.OrganizationDomain
	if err := h.db.Where("domain = ?", host).First(&domain).Error; err == nil {
		return domain.OrganizationID
	}

	// Fall back to global organization
	var globalOrg models.Organization
	if err := h.db.Where("is_global = ?", true).First(&globalOrg).Error; err == nil {
		return globalOrg.ID
	}

	return 0
}

// Redirect handles short URL redirects
// Resolves the organization from the Host header, then looks up the link by (org_id, slug).
// Public links redirect without authentication.
// Private links also redirect (the URL itself is not secret, just the metadata).
// Click count is incremented for all redirects.
func (h *Handler) Redirect(c *gin.Context) {
	slug := c.Param("slug")

	// Resolve organization from Host header
	orgID := h.resolveOrgFromHost(c)
	if orgID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Organization not found"})
		return
	}

	// Find the link within the resolved organization
	var link models.Link
	if err := h.db.Where("organization_id = ? AND slug = ?", orgID, slug).First(&link).Error; err != nil {
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
