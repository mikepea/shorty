package models

import (
	"time"

	"gorm.io/gorm"
)

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	KeyHash     string         `gorm:"not null" json:"-"`
	KeyPrefix   string         `gorm:"not null" json:"key_prefix"` // First few chars for identification
	Description string         `json:"description"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedByID uint           `gorm:"not null" json:"created_by_id"`

	// Relationships
	User      User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedBy User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}
