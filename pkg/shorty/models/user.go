package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemRole represents a user's system-wide role
type SystemRole string

const (
	SystemRoleAdmin SystemRole = "admin"
	SystemRoleUser  SystemRole = "user"
)

// User represents a user in the system
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Name         string         `gorm:"not null" json:"name"`
	SystemRole   SystemRole     `gorm:"type:varchar(20);default:'user'" json:"system_role"`

	// Relationships
	GroupMemberships []GroupMembership `gorm:"foreignKey:UserID" json:"group_memberships,omitempty"`
	APIKeys          []APIKey          `gorm:"foreignKey:UserID" json:"api_keys,omitempty"`
}
