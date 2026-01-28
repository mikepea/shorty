package models

import (
	"time"

	"gorm.io/gorm"
)

// Group represents a group that owns links
// Users can belong to multiple groups, and each user has a personal group
// Groups belong to an organization for multi-tenancy scoping
type Group struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID uint           `gorm:"not null;index" json:"organization_id"` // FK to Organization
	ExternalID     string         `gorm:"index" json:"external_id,omitempty"`    // SCIM externalId
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`

	// Relationships
	Organization Organization      `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Members      []GroupMembership `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Links        []Link            `gorm:"foreignKey:GroupID" json:"links,omitempty"`
}
