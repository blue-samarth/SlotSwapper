package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/models"
)

func TestEventCreate_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")

	event := &models.Event{
		Title:     "Test Meeting",
		UserID:    user.ID,
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Status:    models.StatusBusy,
	}

	err := db.Create(event).Error
	assert.NoError(t, err)
	assert.NotZero(t, event.ID)
}

func TestEventBeforeSave_ValidTimeRange(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")

	event := &models.Event{
		Title:     "Test Meeting",
		UserID:    user.ID,
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Status:    models.StatusBusy,
	}

	err := db.Create(event).Error
	assert.NoError(t, err)
}

func TestEventBeforeSave_EndTimeBeforeStartTime(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")

	event := &models.Event{
		Title:     "Test Meeting",
		UserID:    user.ID,
		StartTime: time.Now().Add(2 * time.Hour),
		EndTime:   time.Now().Add(1 * time.Hour), // End before start
		Status:    models.StatusBusy,
	}

	err := db.Create(event).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_time must be after start_time")
}

func TestEventBeforeSave_StartTimeInPast(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")

	event := &models.Event{
		Title:     "Test Meeting",
		UserID:    user.ID,
		StartTime: time.Now().Add(-1 * time.Hour), // In the past
		EndTime:   time.Now().Add(1 * time.Hour),
		Status:    models.StatusBusy,
	}

	err := db.Create(event).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start_time must be in the future")
}

func TestEventIsSwappable_True(t *testing.T) {
	event := &models.Event{
		Status: models.StatusSwappable,
	}

	assert.True(t, event.IsSwappable())
}

func TestEventIsSwappable_False(t *testing.T) {
	event := &models.Event{
		Status: models.StatusBusy,
	}

	assert.False(t, event.IsSwappable())
}

func TestEventCanChangeStatus_BusyToSwappable(t *testing.T) {
	event := &models.Event{
		Status: models.StatusBusy,
	}

	err := event.CanChangeStatus(models.StatusSwappable)
	assert.NoError(t, err)
}

func TestEventCanChangeStatus_SwappableToBusy(t *testing.T) {
	event := &models.Event{
		Status: models.StatusSwappable,
	}

	err := event.CanChangeStatus(models.StatusBusy)
	assert.NoError(t, err)
}

func TestEventCanChangeStatus_CannotSetSwapPending(t *testing.T) {
	event := &models.Event{
		Status: models.StatusBusy,
	}

	err := event.CanChangeStatus(models.StatusSwapPending)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot set status to SWAP_PENDING directly")
}

func TestEventCanChangeStatus_CannotChangeFromSwapPending(t *testing.T) {
	event := &models.Event{
		Status: models.StatusSwapPending,
	}

	err := event.CanChangeStatus(models.StatusBusy)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot change status from SWAP_PENDING")
}

func TestEventSoftDelete(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")
	event := CreateTestEvent(t, db, "Test Event", user.ID, models.StatusBusy)

	// Soft delete
	event.IsDeleted = true
	db.Save(event)

	// Query without soft delete filter - should still find it
	var foundEvent models.Event
	err := db.Unscoped().First(&foundEvent, event.ID).Error
	assert.NoError(t, err)
	assert.True(t, foundEvent.IsDeleted)

	// Query with soft delete filter - should not find it
	var notFoundEvent models.Event
	err = db.Where("id = ? AND is_deleted = false", event.ID).First(&notFoundEvent).Error
	assert.Error(t, err)
}

func TestEventUserRelationship(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user := CreateTestUser(t, db, "testuser", "test@example.com", "password")
	event := CreateTestEvent(t, db, "Test Event", user.ID, models.StatusBusy)

	// Load with user preload
	var loadedEvent models.Event
	err := db.Preload("User").First(&loadedEvent, event.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, loadedEvent.User.ID)
	assert.Equal(t, user.Email, loadedEvent.User.Email)
}

func TestEventStatusEnum(t *testing.T) {
	assert.Equal(t, models.EventStatus("BUSY"), models.StatusBusy)
	assert.Equal(t, models.EventStatus("SWAPPABLE"), models.StatusSwappable)
	assert.Equal(t, models.EventStatus("SWAP_PENDING"), models.StatusSwapPending)
}
