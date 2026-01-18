package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag represents a tag that can be applied to links
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`

	// Relationships
	Links []Link `gorm:"many2many:link_tags;" json:"links,omitempty"`
}
