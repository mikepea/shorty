package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
	"gorm.io/gorm"
)

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID          string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      string         `gorm:"not null;index;type:varchar(24)" json:"user_id"`
	KeyHash     string         `gorm:"not null" json:"-"`
	KeyPrefix   string         `gorm:"not null" json:"key_prefix"` // First few chars for identification
	Description string         `json:"description"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedByID string         `gorm:"not null;type:varchar(24)" json:"created_by_id"`

	// Relationships
	User      User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedBy User `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// BeforeCreate generates a CUID for the API key if not already set
func (ak *APIKey) BeforeCreate(tx *gorm.DB) error {
	if ak.ID == "" {
		ak.ID = cuid.New()
	}
	return nil
}
