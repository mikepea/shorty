package models

import (
	"time"

	"gorm.io/gorm"
)

// OIDCProvider represents an OIDC identity provider configuration
type OIDCProvider struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"uniqueIndex;not null" json:"name"`         // Display name (e.g., "Okta", "Google")
	Slug         string         `gorm:"uniqueIndex;not null" json:"slug"`         // URL-safe identifier
	Issuer       string         `gorm:"not null" json:"issuer"`                   // OIDC issuer URL
	ClientID     string         `gorm:"not null" json:"client_id"`                // OAuth client ID
	ClientSecret string         `gorm:"not null" json:"-"`                        // OAuth client secret (not exposed in JSON)
	Scopes       string         `gorm:"default:'openid profile email'" json:"scopes"` // Space-separated scopes
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	AutoProvision bool          `gorm:"default:true" json:"auto_provision"`       // Auto-create users on first login
}

// OIDCIdentity links a user to an OIDC provider identity
type OIDCIdentity struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	ProviderID uint           `gorm:"not null;index" json:"provider_id"`
	Subject    string         `gorm:"not null" json:"subject"`                  // OIDC subject (sub claim)
	Email      string         `json:"email"`                                     // Email from OIDC (for reference)

	// Relationships
	User     User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider OIDCProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}

// SCIMToken represents a bearer token for SCIM API access
type SCIMToken struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	TokenHash   string         `gorm:"uniqueIndex;not null" json:"-"`            // SHA-256 hash of token
	TokenPrefix string         `gorm:"not null" json:"token_prefix"`             // First 8 chars for identification
	Description string         `json:"description"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
}
