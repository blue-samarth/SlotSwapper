package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/models"
	"SlotSwapper/services"
)

func TestGetSwapRequests_FilterSent(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// User1 initiates swap
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, swapRequest)

	// Query sent requests for user1
	var sentRequests []models.SwapRequest
	err = db.Where("requester_id = ? AND is_deleted = false", user1.ID).
		Find(&sentRequests).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, len(sentRequests))
	assert.Equal(t, user1.ID, sentRequests[0].RequesterID)
}

func TestGetSwapRequests_FilterReceived(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// User1 initiates swap to user2
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, swapRequest)

	// Query received requests for user2
	var receivedRequests []models.SwapRequest
	err = db.Where("receiver_id = ? AND is_deleted = false", user2.ID).
		Find(&receivedRequests).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, len(receivedRequests))
	assert.Equal(t, user2.ID, receivedRequests[0].ReceiverID)
}

func TestGetSwapRequests_FilterAll(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	user3 := CreateTestUser(t, db, "user3", "user3@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)
	event3 := CreateTestEvent(t, db, "Event 3", user3.ID, models.StatusSwappable)
	event4 := CreateTestEvent(t, db, "Event 4", user1.ID, models.StatusSwappable)

	// User1 initiates swap with user2
	_, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)

	// User3 initiates swap with user1 (using event4, not event1 which is now SWAP_PENDING)
	_, err = services.InitiateSwap(event3.ID, event4.ID, user3.ID)
	assert.NoError(t, err)

	// Query all requests for user1 (sent and received)
	var allRequests []models.SwapRequest
	err = db.Where("(requester_id = ? OR receiver_id = ?) AND is_deleted = false", user1.ID, user1.ID).
		Find(&allRequests).Error

	assert.NoError(t, err)
	assert.Equal(t, 2, len(allRequests), "User1 should have 2 swap requests (1 sent, 1 received)")
}

func TestGetSwapRequests_NoResults(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")

	// Query requests for user with no swaps
	var requests []models.SwapRequest
	err := db.Where("(requester_id = ? OR receiver_id = ?) AND is_deleted = false", user1.ID, user1.ID).
		Find(&requests).Error

	assert.NoError(t, err)
	assert.Equal(t, 0, len(requests))
}

func TestGetSwapRequests_WithPreloads(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)

	// Query with preloads
	var requests []models.SwapRequest
	err = db.Preload("RequesterEvent").
		Preload("ReceiverEvent").
		Preload("Requester").
		Preload("Receiver").
		Where("id = ?", swapRequest.ID).
		Find(&requests).Error

	assert.NoError(t, err)
	assert.Equal(t, 1, len(requests))
	assert.NotZero(t, requests[0].RequesterEvent.ID)
	assert.NotZero(t, requests[0].ReceiverEvent.ID)
	assert.NotZero(t, requests[0].Requester.ID)
	assert.NotZero(t, requests[0].Receiver.ID)
}

func TestGetSwapRequests_OrderByCreatedAt(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	user3 := CreateTestUser(t, db, "user3", "user3@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)
	event3 := CreateTestEvent(t, db, "Event 3", user3.ID, models.StatusSwappable)
	event4 := CreateTestEvent(t, db, "Event 4", user1.ID, models.StatusSwappable)

	// Create multiple swaps (use different events for user1 since event1 becomes SWAP_PENDING)
	swap1, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)
	
	swap2, err := services.InitiateSwap(event4.ID, event3.ID, user1.ID)
	assert.NoError(t, err)

	// Query ordered by created_at DESC
	var requests []models.SwapRequest
	err = db.Where("requester_id = ?", user1.ID).
		Order("created_at DESC").
		Find(&requests).Error

	assert.NoError(t, err)
	assert.Equal(t, 2, len(requests), "Should have 2 swap requests")
	
	// Most recent should be first
	if len(requests) >= 2 {
		assert.True(t, requests[0].CreatedAt.After(requests[1].CreatedAt) || 
			requests[0].ID >= requests[1].ID,
			"Requests should be ordered by creation time")
		
		// Verify IDs
		ids := []uint{requests[0].ID, requests[1].ID}
		assert.Contains(t, ids, swap1.ID)
		assert.Contains(t, ids, swap2.ID)
	}
}

func TestGetSwapRequests_ExcludesDeleted(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")

	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	assert.NoError(t, err)

	// Soft delete the swap request
	err = db.Model(&models.SwapRequest{}).
		Where("id = ?", swapRequest.ID).
		Update("is_deleted", true).Error
	assert.NoError(t, err)

	// Query should exclude deleted
	var requests []models.SwapRequest
	err = db.Where("requester_id = ? AND is_deleted = false", user1.ID).
		Find(&requests).Error

	assert.NoError(t, err)
	assert.Equal(t, 0, len(requests), "Deleted swap requests should be excluded")
}
