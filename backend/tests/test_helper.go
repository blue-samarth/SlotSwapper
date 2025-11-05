package tests

import (
	"os"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"SlotSwapper/config"
	"SlotSwapper/models"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	// Set test environment variables
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-chars-long")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet during tests
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&models.User{}, &models.Event{}, &models.SwapRequest{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the config package's DB to use test database
	config.SetDB(db)

	return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// CreateTestUser creates a user for testing
func CreateTestUser(t *testing.T, db *gorm.DB, username, email, password string) *models.User {
	user := &models.User{
		Username: username,
		Email:    email,
	}
	if err := user.HashPassword(password); err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// CreateTestEvent creates an event for testing
func CreateTestEvent(t *testing.T, db *gorm.DB, title string, userID uint, status models.EventStatus) *models.Event {
	event := &models.Event{
		Title:     title,
		UserID:    userID,
		StartTime: parseTime("2025-11-07T10:00:00Z"),
		EndTime:   parseTime("2025-11-07T11:00:00Z"),
		Status:    status,
	}
	if err := db.Create(event).Error; err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	return event
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
