package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
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
	ID        string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    string         `gorm:"not null;uniqueIndex:idx_user_group;type:varchar(24)" json:"user_id"`
	GroupID   string         `gorm:"not null;uniqueIndex:idx_user_group;type:varchar(24)" json:"group_id"`
	Role      GroupRole      `gorm:"type:varchar(20);default:'member'" json:"role"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Group Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

// BeforeCreate generates a CUID for the group membership if not already set
func (gm *GroupMembership) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == "" {
		gm.ID = cuid.New()
	}
	return nil
}
