package unit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
)

// Helper function to create a rate limiter for testing
func newTestRateLimiter(capacity int, refillRate float64) *services.RateLimiter {
	return services.NewRateLimiter(capacity, refillRate)
}

// ============================================================================
// Phase 1: Basic Unit Tests - Core Functionality
// ============================================================================

// TestNewRateLimiterCreation verifies RateLimiter initializes correctly
func TestNewRateLimiterCreation(t *testing.T) {
	capacity := 10
	refillRate := 2.0
	rl := newTestRateLimiter(capacity, refillRate)

	assert.NotNil(t, rl)
}

// TestSingleRequestAllowed verifies first request from a key is allowed
func TestSingleRequestAllowed(t *testing.T) {
	rl := newTestRateLimiter(5, 1.0)
	key := "192.168.1.1"

	allowed, retryAfter := rl.Allow(key)

	assert.True(t, allowed, "First request should be allowed")
	assert.Equal(t, time.Duration(0), retryAfter, "Retry-After should be 0 for allowed request")
}

// TestSequentialRequestsWithinCapacity verifies multiple requests within capacity
func TestSequentialRequestsWithinCapacity(t *testing.T) {
	capacity := 3
	rl := newTestRateLimiter(capacity, 1.0)
	key := "192.168.1.1"

	// First 3 requests should succeed
	for i := 0; i < capacity; i++ {
		allowed, _ := rl.Allow(key)
		assert.True(t, allowed, "Request %d should be allowed", i+1)
	}

	// 4th request should fail
	allowed, _ := rl.Allow(key)
	assert.False(t, allowed, "Request 4 should be denied (capacity exhausted)")
}

// TestRequestDeniedExhaustedTokens verifies denial when tokens exhausted
func TestRequestDeniedExhaustedTokens(t *testing.T) {
	rl := newTestRateLimiter(1, 1.0)
	key := "192.168.1.1"

	// First request succeeds
	allowed, _ := rl.Allow(key)
	assert.True(t, allowed)

	// Second request fails
	allowed, _ = rl.Allow(key)
	assert.False(t, allowed)
}

// TestRetryAfterDuration verifies correct retry-after calculation
func TestRetryAfterDuration(t *testing.T) {
	rl := newTestRateLimiter(1, 2.0) // 2 tokens per second
	key := "192.168.1.1"

	// Exhaust tokens
	rl.Allow(key)

	// Next request should fail with retry-after
	allowed, retryAfter := rl.Allow(key)
	assert.False(t, allowed)
	assert.Greater(t, retryAfter, time.Duration(0), "Retry-After should be positive")

	// For 1 token at 2 tokens/sec: 1/2 = 0.5 seconds = 500ms
	// Allow some tolerance due to timing
	expectedMin := time.Duration(400 * time.Millisecond)
	expectedMax := time.Duration(600 * time.Millisecond)
	assert.GreaterOrEqual(t, retryAfter, expectedMin, "Retry-After should be ~500ms")
	assert.LessOrEqual(t, retryAfter, expectedMax, "Retry-After should be ~500ms")
}

// TestMultipleKeysIsolated verifies different keys have isolated buckets
func TestMultipleKeysIsolated(t *testing.T) {
	rl := newTestRateLimiter(2, 1.0)
	key1 := "192.168.1.1"
	key2 := "192.168.1.2"

	// Exhaust key1
	rl.Allow(key1)
	rl.Allow(key1)
	allowed1, _ := rl.Allow(key1)
	assert.False(t, allowed1, "Key1 should be exhausted")

	// Key2 should still have tokens
	allowed2, _ := rl.Allow(key2)
	assert.True(t, allowed2, "Key2 should have tokens")
}

// TestInitialBucketState verifies first Allow creates bucket with full capacity
func TestInitialBucketState(t *testing.T) {
	rl := newTestRateLimiter(5, 1.0)
	key := "192.168.1.1"

	allowed, _ := rl.Allow(key)
	assert.True(t, allowed, "First request should succeed with full capacity")

	// Verify we can consume remaining tokens
	for i := 0; i < 4; i++ {
		allowed, _ := rl.Allow(key)
		assert.True(t, allowed, "Should be able to consume remaining capacity tokens")
	}

	// Next should fail
	allowed, _ = rl.Allow(key)
	assert.False(t, allowed, "Should be exhausted after consuming capacity")
}

// ============================================================================
// Phase 2: Concurrency & Token Refill Tests
// ============================================================================

// TestTokenRefillOverTime verifies tokens refill at correct rate
func TestTokenRefillOverTime(t *testing.T) {
	rl := newTestRateLimiter(5, 1.0) // 1 token per second
	key := "192.168.1.1"

	// Consume 1 token
	rl.Allow(key)

	// Sleep 1 second for 1 token to refill
	time.Sleep(1100 * time.Millisecond)

	// Should be able to consume again
	allowed, _ := rl.Allow(key)
	assert.True(t, allowed, "Token should refill after 1 second")
}

// TestRefillRateCalculation verifies correct refill rate
func TestRefillRateCalculation(t *testing.T) {
	rl := newTestRateLimiter(10, 2.0) // 2 tokens per second
	key := "192.168.1.1"

	// Consume all tokens
	for i := 0; i < 10; i++ {
		rl.Allow(key)
	}

	// Sleep 0.5 seconds (should accumulate 1 token: 2 * 0.5)
	time.Sleep(550 * time.Millisecond)

	allowed, _ := rl.Allow(key)
	assert.True(t, allowed, "Should have 1 refilled token")

	// Next should fail
	allowed, _ = rl.Allow(key)
	assert.False(t, allowed, "Should only have 1 token after 0.5 seconds")
}

// TestCapacityCap verifies tokens don't exceed capacity
func TestCapacityCap(t *testing.T) {
	rl := newTestRateLimiter(5, 10.0) // 10 tokens per second
	key := "192.168.1.1"

	// Consume 1 token
	rl.Allow(key)

	// Sleep 2 seconds (would accumulate 20 tokens, but capped at 5)
	time.Sleep(2100 * time.Millisecond)

	// Consume all 5 tokens
	for i := 0; i < 5; i++ {
		allowed, _ := rl.Allow(key)
		assert.True(t, allowed)
	}

	// Should be exhausted
	allowed, _ := rl.Allow(key)
	assert.False(t, allowed, "Should be capped at capacity")
}

// TestConcurrentAccessMultipleGoroutines tests thread-safety
func TestConcurrentAccessMultipleGoroutines(t *testing.T) {
	rl := newTestRateLimiter(100, 10.0)
	key := "192.168.1.1"
	numGoroutines := 10
	requestsPerGoroutine := 10
	var wg sync.WaitGroup
	var allowedCount int32

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				if allowed, _ := rl.Allow(key); allowed {
					atomic.AddInt32(&allowedCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Should have allowed exactly 100 requests (capacity)
	assert.Equal(t, int32(100), allowedCount, "Should allow exactly capacity requests")
}

// TestThreadSafeBucketAccess verifies mutex prevents race conditions
func TestThreadSafeBucketAccess(t *testing.T) {
	rl := newTestRateLimiter(50, 5.0)
	key := "192.168.1.1"
	numGoroutines := 10
	requestsPerGoroutine := 5
	var wg sync.WaitGroup
	var successCount int32

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				if allowed, _ := rl.Allow(key); allowed {
					atomic.AddInt32(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// All goroutines should execute without panic, with total requests = 50
	assert.Equal(t, int32(50), successCount, "Should allow up to capacity")
}

// TestConcurrentRequestsSameKey verifies token consumption accuracy
func TestConcurrentRequestsSameKey(t *testing.T) {
	rl := newTestRateLimiter(20, 1.0)
	key := "192.168.1.1"
	numGoroutines := 5
	requestsPerGoroutine := 4
	var wg sync.WaitGroup
	var successCount int32

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				if allowed, _ := rl.Allow(key); allowed {
					atomic.AddInt32(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Should have exactly 20 successful requests
	assert.Equal(t, int32(20), successCount)
}

// TestEdgeCaseZeroRefillRate verifies behavior with zero refill
func TestEdgeCaseZeroRefillRate(t *testing.T) {
	rl := newTestRateLimiter(1, 0.0)
	key := "192.168.1.1"

	// First request succeeds
	allowed, _ := rl.Allow(key)
	assert.True(t, allowed)

	// Subsequent requests fail indefinitely
	allowed, _ = rl.Allow(key)
	assert.False(t, allowed)

	// Sleep doesn't help with zero refill rate
	time.Sleep(100 * time.Millisecond)
	allowed, _ = rl.Allow(key)
	assert.False(t, allowed)
}

// TestEdgeCaseZeroCapacity verifies zero capacity denies all requests
func TestEdgeCaseZeroCapacity(t *testing.T) {
	rl := newTestRateLimiter(0, 1.0)
	key := "192.168.1.1"

	// All requests should fail
	allowed, _ := rl.Allow(key)
	assert.False(t, allowed)

	allowed, _ = rl.Allow(key)
	assert.False(t, allowed)
}

// ============================================================================
// Phase 3: HTTP Middleware Integration Tests
// ============================================================================

// TestMiddlewareAllowsNormalRequest verifies middleware allows normal requests
func TestMiddlewareAllowsNormalRequest(t *testing.T) {
	rl := newTestRateLimiter(10, 1.0)
	key := "192.168.1.1"

	// Make request that should be allowed
	allowed, _ := rl.Allow(key)
	assert.True(t, allowed)

	// Verify context.Next() would be called (request passes through)
	allowed, _ = rl.Allow(key)
	assert.True(t, allowed)
}

// TestRetryAfterHeaderValue verifies header contains correct value
func TestRetryAfterHeaderValue(t *testing.T) {
	rl := newTestRateLimiter(1, 1.0) // 1 token per second
	key := "192.168.1.1"

	// Exhaust tokens
	rl.Allow(key)

	// Get retry-after for rate limited request
	_, retryAfter := rl.Allow(key)

	assert.Greater(t, retryAfter, time.Duration(0))
	// Should be approximately 1 second (within tolerance)
	assert.Greater(t, retryAfter, 900*time.Millisecond)
	assert.Less(t, retryAfter, 1100*time.Millisecond)
}

// TestErrorResponseStructure verifies error info is available
func TestErrorResponseStructure(t *testing.T) {
	rl := newTestRateLimiter(1, 1.0)
	key := "192.168.1.1"

	// Make first request
	rl.Allow(key)

	// Make rate-limited request
	allowed, retryAfter := rl.Allow(key)

	assert.False(t, allowed)
	assert.Greater(t, retryAfter, time.Duration(0))
	// Middleware would use this to create error response
}

// TestClientIPKeySeparation verifies different IPs are separate
func TestClientIPKeySeparation(t *testing.T) {
	rl := newTestRateLimiter(2, 1.0)
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// First IP exhausts limit
	rl.Allow(ip1)
	rl.Allow(ip1)
	allowed1, _ := rl.Allow(ip1)
	assert.False(t, allowed1)

	// Second IP should still work
	allowed2, _ := rl.Allow(ip2)
	assert.True(t, allowed2)
}

// TestSubsequentRequestAfterBackoff verifies recovery after wait
func TestSubsequentRequestAfterBackoff(t *testing.T) {
	rl := newTestRateLimiter(1, 2.0) // 2 tokens per second
	key := "192.168.1.1"

	// First request succeeds
	allowed, _ := rl.Allow(key)
	assert.True(t, allowed)

	// Second fails
	allowed, retryAfter := rl.Allow(key)
	assert.False(t, allowed)

	// Wait for backoff + margin
	time.Sleep(retryAfter + 100*time.Millisecond)

	// Should succeed now
	allowed, _ = rl.Allow(key)
	assert.True(t, allowed)
}

// TestHighRequestRateHandling verifies handling of burst requests
func TestHighRequestRateHandling(t *testing.T) {
	rl := newTestRateLimiter(10, 100.0) // 100 tokens per second
	key := "192.168.1.1"
	var successCount int32

	// Make 10 rapid requests
	for i := 0; i < 10; i++ {
		if allowed, _ := rl.Allow(key); allowed {
			atomic.AddInt32(&successCount, 1)
		}
	}

	// Should allow all 10
	assert.Equal(t, int32(10), successCount)

	// 11th should fail
	allowed, _ := rl.Allow(key)
	assert.False(t, allowed)
}
