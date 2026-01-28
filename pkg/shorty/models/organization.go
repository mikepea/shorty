package models

import (
	"time"

	"gorm.io/gorm"
)

// OrgRole represents a user's role within an organization
type OrgRole string

const (
	OrgRoleAdmin  OrgRole = "admin"
	OrgRoleMember OrgRole = "member"
)

// Organization represents a tenant in the multi-tenancy system.
// Organizations scope SSO settings, SCIM provisioning, teams/groups, and link slugs.
// There is always a special "Shorty Global" organization (IsGlobal=true) that serves
// as the default for public signups and unrecognized domains.
type Organization struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"not null" json:"name"`            // Display name (e.g., "Acme Corp")
	Slug      string         `gorm:"uniqueIndex;not null" json:"slug"` // URL-safe identifier, unique across all orgs
	IsGlobal  bool           `gorm:"default:false" json:"is_global"`   // True only for "Shorty Global"

	// Relationships
	Members []OrganizationMembership `gorm:"foreignKey:OrganizationID" json:"members,omitempty"`
	Domains []OrganizationDomain     `gorm:"foreignKey:OrganizationID" json:"domains,omitempty"`
	Groups  []Group                  `gorm:"foreignKey:OrganizationID" json:"groups,omitempty"`
}

// OrganizationMembership represents the many-to-many relationship between users and organizations.
// Users can belong to multiple organizations with different roles in each.
type OrganizationMembership struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID uint           `gorm:"not null;uniqueIndex:idx_org_user" json:"organization_id"`
	UserID         uint           `gorm:"not null;uniqueIndex:idx_org_user" json:"user_id"`
	Role           OrgRole        `gorm:"type:varchar(20);default:'member'" json:"role"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// OrganizationDomain represents a domain that maps to an organization.
// When a request comes in, the Host header is matched against these domains
// to determine which organization's links to serve.
// Multiple domains can map to the same organization.
type OrganizationDomain struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID uint           `gorm:"not null;index" json:"organization_id"`
	Domain         string         `gorm:"uniqueIndex;not null" json:"domain"` // e.g., "go.acme.com" - unique across all orgs
	IsPrimary      bool           `gorm:"default:false" json:"is_primary"`    // Primary domain for generating URLs

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}
