package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
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
	ID        string         `gorm:"primarykey;type:varchar(24)" json:"id"`
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

// BeforeCreate generates a CUID for the organization if not already set
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = cuid.New()
	}
	return nil
}

// OrganizationMembership represents the many-to-many relationship between users and organizations.
// Users can belong to multiple organizations with different roles in each.
type OrganizationMembership struct {
	ID             string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID string         `gorm:"not null;uniqueIndex:idx_org_user;type:varchar(24)" json:"organization_id"`
	UserID         string         `gorm:"not null;uniqueIndex:idx_org_user;type:varchar(24)" json:"user_id"`
	Role           OrgRole        `gorm:"type:varchar(20);default:'member'" json:"role"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate generates a CUID for the organization membership if not already set
func (om *OrganizationMembership) BeforeCreate(tx *gorm.DB) error {
	if om.ID == "" {
		om.ID = cuid.New()
	}
	return nil
}

// OrganizationDomain represents a domain that maps to an organization.
// When a request comes in, the Host header is matched against these domains
// to determine which organization's links to serve.
// Multiple domains can map to the same organization.
type OrganizationDomain struct {
	ID             string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID string         `gorm:"not null;index;type:varchar(24)" json:"organization_id"`
	Domain         string         `gorm:"uniqueIndex;not null" json:"domain"` // e.g., "go.acme.com" - unique across all orgs
	IsPrimary      bool           `gorm:"default:false" json:"is_primary"`    // Primary domain for generating URLs

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// BeforeCreate generates a CUID for the organization domain if not already set
func (od *OrganizationDomain) BeforeCreate(tx *gorm.DB) error {
	if od.ID == "" {
		od.ID = cuid.New()
	}
	return nil
}
