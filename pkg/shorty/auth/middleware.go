package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

const (
	// ContextKeyUserID is the key for user ID in gin context
	ContextKeyUserID = "user_id"
	// ContextKeyEmail is the key for email in gin context
	ContextKeyEmail = "email"
	// ContextKeySystemRole is the key for system role in gin context
	ContextKeySystemRole = "system_role"
	// ContextKeyOrgID is the key for organization ID in gin context
	ContextKeyOrgID = "organization_id"
	// ContextKeyOrgRole is the key for organization role in gin context
	ContextKeyOrgRole = "organization_role"
)

// AuthMiddleware validates JWT tokens and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expect "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ValidateToken(tokenString)
		if err != nil {
			if err == ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeySystemRole, claims.SystemRole)

		c.Next()
	}
}

// RequireAdmin middleware checks if the user has admin system role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeySystemRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID returns the user ID from the gin context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetEmail returns the email from the gin context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(ContextKeyEmail)
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetSystemRole returns the system role from the gin context
func GetSystemRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeySystemRole)
	if !exists {
		return "", false
	}
	return role.(string), true
}

// GetOrgID returns the organization ID from the gin context
func GetOrgID(c *gin.Context) (uint, bool) {
	orgID, exists := c.Get(ContextKeyOrgID)
	if !exists {
		return 0, false
	}
	return orgID.(uint), true
}

// GetOrgRole returns the organization role from the gin context
func GetOrgRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeyOrgRole)
	if !exists {
		return "", false
	}
	return role.(string), true
}

// OrgMiddleware validates the X-Organization-ID header and sets org context.
// The user must be a member of the specified organization.
// If no header is provided, it defaults to the global organization.
func OrgMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get organization ID from header or query param
		orgIDStr := c.GetHeader("X-Organization-ID")
		if orgIDStr == "" {
			orgIDStr = c.Query("org_id")
		}

		var orgID uint
		var membership models.OrganizationMembership

		if orgIDStr != "" {
			// Parse the provided org ID
			parsed, err := strconv.ParseUint(orgIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
				c.Abort()
				return
			}
			orgID = uint(parsed)

			// Verify membership
			if err := db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&membership).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this organization"})
				c.Abort()
				return
			}
		} else {
			// Default to global organization
			var globalOrg models.Organization
			if err := db.Where("is_global = ?", true).First(&globalOrg).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Global organization not found"})
				c.Abort()
				return
			}
			orgID = globalOrg.ID

			// Check if user is a member of global org, create membership if not
			if err := db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&membership).Error; err != nil {
				// User not yet a member of global org - auto-add them as member
				membership = models.OrganizationMembership{
					OrganizationID: orgID,
					UserID:         userID.(uint),
					Role:           models.OrgRoleMember,
				}
				if err := db.Create(&membership).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to global organization"})
					c.Abort()
					return
				}
			}
		}

		// Set organization context
		c.Set(ContextKeyOrgID, orgID)
		c.Set(ContextKeyOrgRole, string(membership.Role))

		c.Next()
	}
}

// RequireOrgAdmin middleware checks if the user is an admin of the current organization
func RequireOrgAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyOrgRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization context required"})
			c.Abort()
			return
		}

		if role != string(models.OrgRoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
