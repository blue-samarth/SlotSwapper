package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/models"
	"SlotSwapper/utils"
)

func TestUserHashPassword(t *testing.T) {
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	password := "password123"
	err := user.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, password, user.Password) // Should be hashed
}

func TestUserCheckPassword_Success(t *testing.T) {
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	password := "password123"
	user.HashPassword(password)

	result := user.CheckPassword(password)
	assert.True(t, result)
}

func TestUserCheckPassword_Failure(t *testing.T) {
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	user.HashPassword("password123")

	result := user.CheckPassword("wrongpassword")
	assert.False(t, result)
}

func TestUserBeforeSave_NormalizesEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := &models.User{
		Username: "testuser",
		Email:    "TEST@EXAMPLE.COM", // Uppercase
		Password: "hashedpassword",
	}

	err := db.Create(user).Error
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email) // Should be lowercase
}

func TestUserBeforeSave_NormalizesUsername(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := &models.User{
		Username: "TestUser", // Mixed case
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	err := db.Create(user).Error
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username) // Should be lowercase
}

func TestUserUniqueEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "test@example.com", "password123")
	assert.NotNil(t, user1)

	// Try to create another user with same email
	user2 := &models.User{
		Username: "user2",
		Email:    "test@example.com", // Duplicate email
	}
	user2.HashPassword("password456")

	err := db.Create(user2).Error
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestJWTGenerateAndValidate(t *testing.T) {
	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	// Generate token
	token, err := utils.GenerateJWTToken(userID, email, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := utils.ValidateJWTToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, username, claims.Username)
}

func TestJWTValidate_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"

	claims, err := utils.ValidateJWTToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTValidate_ExpiredToken(t *testing.T) {
	// This would require mocking time, skipping for now
	// In real scenario, you'd use a library like github.com/benbjohnson/clock
	t.Skip("Requires time mocking for expired token test")
}

func TestGetJWTSecret_TooShort(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	os.Setenv("JWT_SECRET", "short") // Less than 32 chars

	_, err := utils.GetJWTSecret()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 32 characters")
}
