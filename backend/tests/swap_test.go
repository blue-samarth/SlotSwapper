package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/models"
	"SlotSwapper/services"
)

func TestInitiateSwap_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	// Create users and events
	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// Initiate swap
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)

	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}
	if swapRequest == nil {
		t.Fatal("swapRequest is nil but no error returned")
	}
	assert.Equal(t, models.StatusPending, swapRequest.Status)
	assert.Equal(t, user1.ID, swapRequest.RequesterID)
	assert.Equal(t, user2.ID, swapRequest.ReceiverID)

	// Verify events are now SWAP_PENDING
	var updatedEvent1, updatedEvent2 models.Event
	db.First(&updatedEvent1, event1.ID)
	db.First(&updatedEvent2, event2.ID)
	
	assert.Equal(t, models.StatusSwapPending, updatedEvent1.Status)
	assert.Equal(t, models.StatusSwapPending, updatedEvent2.Status)
}

func TestInitiateSwap_EventNotSwappable(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusBusy) // NOT SWAPPABLE
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// Attempt to initiate swap
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)

	assert.Error(t, err)
	assert.Nil(t, swapRequest)
	assert.Contains(t, err.Error(), "not swappable")
}

func TestInitiateSwap_NotOwner(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	user3 := CreateTestUser(t, db, "user3", "user3@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// User3 tries to initiate swap with User1's event
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user3.ID)

	assert.Error(t, err)
	assert.Nil(t, swapRequest)
	assert.Contains(t, err.Error(), "don't own")
}

func TestInitiateSwap_EventNotFound(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")

	// Try to initiate swap with non-existent events
	swapRequest, err := services.InitiateSwap(999, 1000, user1.ID)

	assert.Error(t, err)
	assert.Nil(t, swapRequest)
}

func TestAcceptSwap_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// Initiate swap
	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// Accept swap
	err = services.AcceptSwap(swapRequest.ID, user2.ID)

	assert.NoError(t, err)

	// Verify ownership swap
	var updatedEvent1, updatedEvent2 models.Event
	db.First(&updatedEvent1, event1.ID)
	db.First(&updatedEvent2, event2.ID)
	
	assert.Equal(t, user2.ID, updatedEvent1.UserID) // Event1 now belongs to User2
	assert.Equal(t, user1.ID, updatedEvent2.UserID) // Event2 now belongs to User1
	assert.Equal(t, models.StatusBusy, updatedEvent1.Status)
	assert.Equal(t, models.StatusBusy, updatedEvent2.Status)

	// Verify swap request status
	var updatedSwap models.SwapRequest
	db.First(&updatedSwap, swapRequest.ID)
	assert.Equal(t, models.StatusAccepted, updatedSwap.Status)
}

func TestAcceptSwap_NotReceiver(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	user3 := CreateTestUser(t, db, "user3", "user3@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// User3 tries to accept (not the receiver)
	err = services.AcceptSwap(swapRequest.ID, user3.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not the receiver")
}

func TestAcceptSwap_AlreadyAccepted(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}
	
	// Accept once
	err = services.AcceptSwap(swapRequest.ID, user2.ID)
	if err != nil {
		t.Fatalf("First AcceptSwap failed: %v", err)
	}

	// Try to accept again
	err = services.AcceptSwap(swapRequest.ID, user2.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not pending")
}

func TestRejectSwap_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// Reject swap
	err = services.RejectSwap(swapRequest.ID, user2.ID)

	assert.NoError(t, err)

	// Verify events reverted to SWAPPABLE
	var updatedEvent1, updatedEvent2 models.Event
	db.First(&updatedEvent1, event1.ID)
	db.First(&updatedEvent2, event2.ID)
	
	assert.Equal(t, models.StatusSwappable, updatedEvent1.Status)
	assert.Equal(t, models.StatusSwappable, updatedEvent2.Status)

	// Verify swap status
	var updatedSwap models.SwapRequest
	db.First(&updatedSwap, swapRequest.ID)
	assert.Equal(t, models.StatusRejected, updatedSwap.Status)
}

func TestCancelSwap_Success(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// Requester cancels
	err = services.CancelSwap(swapRequest.ID, user1.ID)

	assert.NoError(t, err)

	// Verify events reverted to SWAPPABLE
	var updatedEvent1, updatedEvent2 models.Event
	db.First(&updatedEvent1, event1.ID)
	db.First(&updatedEvent2, event2.ID)
	
	assert.Equal(t, models.StatusSwappable, updatedEvent1.Status)
	assert.Equal(t, models.StatusSwappable, updatedEvent2.Status)

	// Verify swap status
	var updatedSwap models.SwapRequest
	db.First(&updatedSwap, swapRequest.ID)
	assert.Equal(t, models.StatusCancelled, updatedSwap.Status)
}

func TestCancelSwap_NotRequester(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// Receiver tries to cancel (only requester can)
	err = services.CancelSwap(swapRequest.ID, user2.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not the requester")
}

// RACE CONDITION TESTS

func TestConcurrentSwapAttempts(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	user3 := CreateTestUser(t, db, "user3", "user3@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)
	event3 := CreateTestEvent(t, db, "Event 3", user3.ID, models.StatusSwappable)

	// Two users try to swap with the same event simultaneously
	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)
	
	go func() {
		defer wg.Done()
		_, err1 = services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	}()
	
	go func() {
		defer wg.Done()
		_, err2 = services.InitiateSwap(event3.ID, event2.ID, user3.ID)
	}()

	wg.Wait()

	// Only one should succeed
	successCount := 0
	if err1 == nil {
		successCount++
	}
	if err2 == nil {
		successCount++
	}

	assert.Equal(t, 1, successCount, "Only one concurrent swap should succeed")
}

func TestConcurrentAcceptAttempts(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	swapRequest, err := services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	if err != nil {
		t.Fatalf("InitiateSwap failed: %v", err)
	}

	// Try to accept the same swap multiple times concurrently
	var wg sync.WaitGroup
	errors := make([]error, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			errors[index] = services.AcceptSwap(swapRequest.ID, user2.ID)
		}(i)
	}

	wg.Wait()

	// Only one should succeed
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}

	assert.Equal(t, 1, successCount, "Only one accept should succeed")
}

func TestSwapRequestStatusTransitions(t *testing.T) {
	swapRequest := &models.SwapRequest{
		Status: models.StatusPending,
	}

	// Can transition from PENDING to ACCEPTED
	assert.True(t, swapRequest.CanChangeStatus(models.StatusAccepted))

	// Can transition from PENDING to REJECTED
	assert.True(t, swapRequest.CanChangeStatus(models.StatusRejected))

	// Can transition from PENDING to CANCELLED
	assert.True(t, swapRequest.CanChangeStatus(models.StatusCancelled))

	// Cannot transition from ACCEPTED
	swapRequest.Status = models.StatusAccepted
	assert.False(t, swapRequest.CanChangeStatus(models.StatusPending))
	assert.False(t, swapRequest.CanChangeStatus(models.StatusRejected))
}

func TestDeadlockPrevention(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	user1 := CreateTestUser(t, db, "user1", "user1@test.com", "password")
	user2 := CreateTestUser(t, db, "user2", "user2@test.com", "password")
	
	event1 := CreateTestEvent(t, db, "Event 1", user1.ID, models.StatusSwappable)
	event2 := CreateTestEvent(t, db, "Event 2", user2.ID, models.StatusSwappable)

	// Try to initiate swaps in opposite order simultaneously
	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)
	
	go func() {
		defer wg.Done()
		_, err1 = services.InitiateSwap(event1.ID, event2.ID, user1.ID)
	}()
	
	go func() {
		defer wg.Done()
		_, err2 = services.InitiateSwap(event2.ID, event1.ID, user2.ID)
	}()

	wg.Wait()

	// Both should complete without deadlock (one succeeds, one fails)
	assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil))
}
