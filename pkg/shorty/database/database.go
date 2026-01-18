package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes the database connection.
// For now, uses SQLite. Can be swapped to Spanner later via GORM driver.
func Connect(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// GetDB returns the database instance.
func GetDB() *gorm.DB {
	return DB
}
