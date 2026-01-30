package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
	"gorm.io/gorm"
)

// Tag represents a tag that can be applied to links
type Tag struct {
	ID        string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`

	// Relationships
	Links []Link `gorm:"many2many:link_tags;" json:"links,omitempty"`
}

// BeforeCreate generates a CUID for the tag if not already set
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = cuid.New()
	}
	return nil
}
