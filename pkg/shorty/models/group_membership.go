package models

import (
	"time"

	"gorm.io/gorm"
)

// GroupRole represents a user's role within a specific group
type GroupRole string

const (
	GroupRoleAdmin  GroupRole = "admin"
	GroupRoleMember GroupRole = "member"
)

// GroupMembership represents the many-to-many relationship between users and groups
type GroupMembership struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null;uniqueIndex:idx_user_group" json:"user_id"`
	GroupID   uint           `gorm:"not null;uniqueIndex:idx_user_group" json:"group_id"`
	Role      GroupRole      `gorm:"type:varchar(20);default:'member'" json:"role"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Group Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}
