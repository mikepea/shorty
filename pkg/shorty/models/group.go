package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
	"gorm.io/gorm"
)

// Group represents a group that owns links
// Users can belong to multiple groups, and each user has a personal group
// Groups belong to an organization for multi-tenancy scoping
type Group struct {
	ID             string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID string         `gorm:"not null;index;type:varchar(24)" json:"organization_id"` // FK to Organization
	ExternalID     string         `gorm:"index" json:"external_id,omitempty"`                     // SCIM externalId
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`

	// Relationships
	Organization Organization      `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Members      []GroupMembership `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Links        []Link            `gorm:"foreignKey:GroupID" json:"links,omitempty"`
}

// BeforeCreate generates a CUID for the group if not already set
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = cuid.New()
	}
	return nil
}
