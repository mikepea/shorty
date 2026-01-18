package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func TestAutoMigrate(t *testing.T) {
	db := setupTestDB(t)

	err := AutoMigrate(db)
	if err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	// Verify tables exist by checking if we can query them
	tables := []string{"users", "groups", "group_memberships", "links", "tags", "api_keys", "link_tags"}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Expected table %s to exist", table)
		}
	}
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)
	AutoMigrate(db)

	user := User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test User",
		SystemRole:   SystemRoleUser,
	}

	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("Failed to create user: %v", result.Error)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set after create")
	}

	// Test unique email constraint
	user2 := User{
		Email:        "test@example.com",
		PasswordHash: "another_hash",
		Name:         "Another User",
	}
	result = db.Create(&user2)
	if result.Error == nil {
		t.Error("Expected error when creating user with duplicate email")
	}
}

func TestGroupAndMembership(t *testing.T) {
	db := setupTestDB(t)
	AutoMigrate(db)

	// Create user
	user := User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test User",
	}
	db.Create(&user)

	// Create group
	group := Group{
		Name:        "Test Group",
		Description: "A test group",
	}
	db.Create(&group)

	// Create membership
	membership := GroupMembership{
		UserID:  user.ID,
		GroupID: group.ID,
		Role:    GroupRoleAdmin,
	}
	result := db.Create(&membership)
	if result.Error != nil {
		t.Fatalf("Failed to create membership: %v", result.Error)
	}

	// Verify relationship
	var loadedUser User
	db.Preload("GroupMemberships").First(&loadedUser, user.ID)
	if len(loadedUser.GroupMemberships) != 1 {
		t.Errorf("Expected 1 membership, got %d", len(loadedUser.GroupMemberships))
	}
}

func TestLinkWithTags(t *testing.T) {
	db := setupTestDB(t)
	AutoMigrate(db)

	// Create user and group
	user := User{Email: "test@example.com", PasswordHash: "hash", Name: "Test"}
	db.Create(&user)
	group := Group{Name: "Test Group"}
	db.Create(&group)

	// Create tags
	tag1 := Tag{Name: "golang"}
	tag2 := Tag{Name: "programming"}
	db.Create(&tag1)
	db.Create(&tag2)

	// Create link with tags
	link := Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "test-link",
		URL:         "https://example.com",
		Title:       "Example Site",
		Tags:        []Tag{tag1, tag2},
	}
	result := db.Create(&link)
	if result.Error != nil {
		t.Fatalf("Failed to create link: %v", result.Error)
	}

	// Verify tags relationship
	var loadedLink Link
	db.Preload("Tags").First(&loadedLink, link.ID)
	if len(loadedLink.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(loadedLink.Tags))
	}
}

func TestSlugUniqueness(t *testing.T) {
	db := setupTestDB(t)
	AutoMigrate(db)

	user := User{Email: "test@example.com", PasswordHash: "hash", Name: "Test"}
	db.Create(&user)
	group := Group{Name: "Test Group"}
	db.Create(&group)

	link1 := Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "unique-slug",
		URL:         "https://example1.com",
	}
	db.Create(&link1)

	link2 := Link{
		GroupID:     group.ID,
		CreatedByID: user.ID,
		Slug:        "unique-slug",
		URL:         "https://example2.com",
	}
	result := db.Create(&link2)
	if result.Error == nil {
		t.Error("Expected error when creating link with duplicate slug")
	}
}
