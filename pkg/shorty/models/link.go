package models

import (
	"time"

	"gorm.io/gorm"
)

// Link represents a shortened URL/bookmark
type Link struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	GroupID     uint           `gorm:"not null;index" json:"group_id"`
	CreatedByID uint           `gorm:"not null" json:"created_by_id"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	URL         string         `gorm:"not null" json:"url"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	IsUnread    bool           `gorm:"default:true" json:"is_unread"`
	ClickCount  uint           `gorm:"default:0" json:"click_count"`

	// Relationships
	Group     Group    `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	CreatedBy User     `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Tags      []Tag    `gorm:"many2many:link_tags;" json:"tags,omitempty"`
}
