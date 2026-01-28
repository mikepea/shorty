package models

import "gorm.io/gorm"

// AllModels returns all models for migration
// Note: Organization must be migrated first as other models depend on it
func AllModels() []interface{} {
	return []interface{}{
		&Organization{},
		&OrganizationMembership{},
		&OrganizationDomain{},
		&User{},
		&Group{},
		&GroupMembership{},
		&Link{},
		&Tag{},
		&APIKey{},
		&OIDCProvider{},
		&OIDCIdentity{},
		&SCIMToken{},
	}
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}
