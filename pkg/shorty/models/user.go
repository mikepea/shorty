package models

import (
	"time"

	"github.com/mikepea/shorty/pkg/shorty/cuid"
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
	ID           string         `gorm:"primarykey;type:varchar(24)" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	ExternalID   string         `gorm:"index" json:"external_id,omitempty"`       // SCIM externalId
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `json:"-"`                                        // Optional for OIDC-only users
	Name         string         `gorm:"not null" json:"name"`
	GivenName    string         `json:"given_name,omitempty"`                     // SCIM givenName
	FamilyName   string         `json:"family_name,omitempty"`                    // SCIM familyName
	Active       bool           `gorm:"default:true" json:"active"`               // SCIM active status
	SystemRole   SystemRole     `gorm:"type:varchar(20);default:'user'" json:"system_role"`

	// Relationships
	OrganizationMemberships []OrganizationMembership `gorm:"foreignKey:UserID" json:"organization_memberships,omitempty"`
	GroupMemberships        []GroupMembership        `gorm:"foreignKey:UserID" json:"group_memberships,omitempty"`
	APIKeys                 []APIKey                 `gorm:"foreignKey:UserID" json:"api_keys,omitempty"`
	OIDCIdentities          []OIDCIdentity           `gorm:"foreignKey:UserID" json:"oidc_identities,omitempty"`
}

// BeforeCreate generates a CUID for the user if not already set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = cuid.New()
	}
	return nil
}
