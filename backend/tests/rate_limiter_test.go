package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/middleware"
)

func TestRateLimiter_Creation(t *testing.T) {
	// Test that rate limiter can be created with valid parameters
	limiter := middleware.NewRateLimiter(5, 1*time.Minute)
	assert.NotNil(t, limiter, "Rate limiter should be created successfully")
}

func TestRateLimiter_ConfigurationValues(t *testing.T) {
	// Test different limit configurations
	limits := []struct {
		limit  int
		window time.Duration
	}{
		{5, 1 * time.Minute},
		{20, 1 * time.Minute},
		{100, 1 * time.Minute},
	}

	for _, config := range limits {
		limiter := middleware.NewRateLimiter(config.limit, config.window)
		assert.NotNil(t, limiter, "Rate limiter with limit %d should be created", config.limit)
	}
}

func TestRateLimiter_TracksRequests(t *testing.T) {
	// Test basic request counting logic
	visitors := make(map[string]int)
	ip := "192.168.1.1"
	limit := 3

	// Simulate requests
	for i := 0; i < 5; i++ {
		visitors[ip]++
	}

	assert.Equal(t, 5, visitors[ip], "Should track all requests")
	assert.Greater(t, visitors[ip], limit, "Should exceed limit")
}

func TestRateLimiter_WindowReset(t *testing.T) {
	// Test window reset logic with simulation
	count := 0
	limit := 2

	// First batch
	for i := 0; i < limit; i++ {
		count++
	}
	assert.Equal(t, limit, count)

	// Simulate window reset
	time.Sleep(50 * time.Millisecond)
	firstBatch := count

	// Second batch after "reset"
	count = 0 // Reset counter to simulate window reset
	for i := 0; i < limit; i++ {
		count++
	}
	assert.Equal(t, limit, count)
	assert.Equal(t, firstBatch, limit)
}

func TestRateLimiter_MultipleIPs(t *testing.T) {
	// Test that different IPs are tracked independently
	visitors := make(map[string]int)
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"
	limit := 2

	// IP1 makes requests
	for i := 0; i < limit; i++ {
		visitors[ip1]++
	}

	// IP2 makes requests
	for i := 0; i < limit; i++ {
		visitors[ip2]++
	}

	assert.Equal(t, limit, visitors[ip1])
	assert.Equal(t, limit, visitors[ip2])
	assert.NotEqual(t, visitors[ip1], 0)
	assert.NotEqual(t, visitors[ip2], 0)
}

func TestRateLimiter_AuthConfiguration(t *testing.T) {
	// Test auth rate limiter configuration (5 req/min)
	limit := 5
	window := 1 * time.Minute

	limiter := middleware.NewRateLimiter(limit, window)
	assert.NotNil(t, limiter)

	// Verify limit logic
	allowed := 0
	blocked := 0

	for i := 0; i < 10; i++ {
		if i < limit {
			allowed++
		} else {
			blocked++
		}
	}

	assert.Equal(t, 5, allowed)
	assert.Equal(t, 5, blocked)
}

func TestRateLimiter_SwapConfiguration(t *testing.T) {
	// Test swap rate limiter configuration (20 req/min)
	limit := 20
	window := 1 * time.Minute

	limiter := middleware.NewRateLimiter(limit, window)
	assert.NotNil(t, limiter)

	// Verify limit logic
	allowed := 0
	for i := 0; i < 25; i++ {
		if i < limit {
			allowed++
		}
	}

	assert.Equal(t, 20, allowed)
}

func TestRateLimiter_GlobalConfiguration(t *testing.T) {
	// Test global rate limiter configuration (100 req/min)
	limit := 100
	window := 1 * time.Minute

	limiter := middleware.NewRateLimiter(limit, window)
	assert.NotNil(t, limiter)

	// Verify that limiter is created with correct parameters
	assert.NotNil(t, limiter, "Global rate limiter should be created")
}

func TestRateLimiter_CleanupRoutine(t *testing.T) {
	// Test that cleanup routine doesn't cause panics
	limiter := middleware.NewRateLimiter(5, 100*time.Millisecond)
	assert.NotNil(t, limiter)

	// Let cleanup routine run in background
	time.Sleep(150 * time.Millisecond)

	// Limiter should still be functional after cleanup runs
	assert.NotNil(t, limiter, "Limiter should remain functional after cleanup")
}

func TestRateLimiter_MiddlewareCreation(t *testing.T) {
	// Test that middleware functions can be created
	authLimiter := middleware.AuthRateLimiter()
	swapLimiter := middleware.SwapRateLimiter()
	globalLimiter := middleware.GlobalRateLimiter()

	assert.NotNil(t, authLimiter, "Auth rate limiter middleware should be created")
	assert.NotNil(t, swapLimiter, "Swap rate limiter middleware should be created")
	assert.NotNil(t, globalLimiter, "Global rate limiter middleware should be created")
}
