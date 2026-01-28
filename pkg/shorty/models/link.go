package models

import (
	"time"

	"gorm.io/gorm"
)

// Link represents a shortened URL/bookmark
// Links are scoped to organizations - the same slug can exist in different organizations
type Link struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	OrganizationID uint           `gorm:"not null;uniqueIndex:idx_org_slug" json:"organization_id"` // FK to Organization (denormalized from Group for direct queries)
	GroupID        uint           `gorm:"not null;index" json:"group_id"`
	CreatedByID    uint           `gorm:"not null" json:"created_by_id"`
	Slug           string         `gorm:"not null;uniqueIndex:idx_org_slug" json:"slug"` // Unique within organization
	URL            string         `gorm:"not null" json:"url"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	IsPublic       bool           `gorm:"default:false" json:"is_public"`
	IsUnread       bool           `gorm:"default:true" json:"is_unread"`
	ClickCount     uint           `gorm:"default:0" json:"click_count"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Group        Group        `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	CreatedBy    User         `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Tags         []Tag        `gorm:"many2many:link_tags;" json:"tags,omitempty"`
}
