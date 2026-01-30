package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
	"gorm.io/gorm"
)

// OIDCProvider represents an OIDC identity provider configuration
// Providers are scoped to organizations - each organization can have its own SSO settings
type OIDCProvider struct {
	ID             string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID string         `gorm:"not null;uniqueIndex:idx_org_provider_name;uniqueIndex:idx_org_provider_slug;type:varchar(24)" json:"organization_id"` // FK to Organization
	Name           string         `gorm:"not null;uniqueIndex:idx_org_provider_name" json:"name"`                                                               // Display name (e.g., "Okta", "Google") - unique per org
	Slug           string         `gorm:"not null;uniqueIndex:idx_org_provider_slug" json:"slug"`                                                               // URL-safe identifier - unique per org
	Issuer         string         `gorm:"not null" json:"issuer"`                                                                                               // OIDC issuer URL
	ClientID       string         `gorm:"not null" json:"client_id"`                                                                                            // OAuth client ID
	ClientSecret   string         `gorm:"not null" json:"-"`                                                                                                    // OAuth client secret (not exposed in JSON)
	Scopes         string         `gorm:"default:'openid profile email'" json:"scopes"`                                                                         // Space-separated scopes
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	AutoProvision  bool           `gorm:"default:true" json:"auto_provision"` // Auto-create users on first login

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// BeforeCreate generates a CUID for the OIDC provider if not already set
func (op *OIDCProvider) BeforeCreate(tx *gorm.DB) error {
	if op.ID == "" {
		op.ID = cuid.New()
	}
	return nil
}

// OIDCIdentity links a user to an OIDC provider identity
type OIDCIdentity struct {
	ID         string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     string         `gorm:"not null;index;type:varchar(24)" json:"user_id"`
	ProviderID string         `gorm:"not null;index;type:varchar(24)" json:"provider_id"`
	Subject    string         `gorm:"not null" json:"subject"` // OIDC subject (sub claim)
	Email      string         `json:"email"`                   // Email from OIDC (for reference)

	// Relationships
	User     User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider OIDCProvider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
}

// BeforeCreate generates a CUID for the OIDC identity if not already set
func (oi *OIDCIdentity) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = cuid.New()
	}
	return nil
}

// SCIMToken represents a bearer token for SCIM API access
// Tokens are scoped to organizations - SCIM provisioning is per-organization
type SCIMToken struct {
	ID             string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID string         `gorm:"not null;index;type:varchar(24)" json:"organization_id"` // FK to Organization
	TokenHash      string         `gorm:"uniqueIndex;not null" json:"-"`                          // SHA-256 hash of token
	TokenPrefix    string         `gorm:"not null" json:"token_prefix"`                           // First 8 chars for identification
	Description    string         `json:"description"`
	LastUsedAt     *time.Time     `json:"last_used_at"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// BeforeCreate generates a CUID for the SCIM token if not already set
func (st *SCIMToken) BeforeCreate(tx *gorm.DB) error {
	if st.ID == "" {
		st.ID = cuid.New()
	}
	return nil
}
