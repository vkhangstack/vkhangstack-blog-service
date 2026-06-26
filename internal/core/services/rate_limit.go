package services

import (
	"math"
	"sync"
	"time"
)

// Implement token bucket algorithm for rate limiting
type bucket struct {
	tokens         float64
	capacity       float64
	refillRate     float64
	lastRefillTime time.Time
}
type RateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*bucket
	capacity   float64
	refillRate float64
	idleTTL    time.Duration
}

func (l *RateLimiter) cleanupIdleBuckets() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-l.idleTTL)

		l.mu.Lock()
		for key, b := range l.buckets {
			if b.lastRefillTime.Before(cutoff) {
				delete(l.buckets, key)
			}
		}
		l.mu.Unlock()

	}
}

func NewRateLimiter(capacity int, refillRate float64) *RateLimiter {
	t := &RateLimiter{
		buckets:    make(map[string]*bucket),
		capacity:   float64(capacity),
		refillRate: refillRate,
		idleTTL:    5 * time.Minute,
	}
	go t.cleanupIdleBuckets()
	return t
}

func (rl *RateLimiter) Allow(key string) (bool, time.Duration) {
	// Implement rate limiting logic here
	// For example, you can check if the number of requests from the given IP exceeds the limit
	// and return true or false accordingly.
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[key]
	if !ok {
		b = &bucket{
			tokens:         rl.capacity,
			capacity:       rl.capacity,
			refillRate:     rl.refillRate,
			lastRefillTime: now,
		}
		rl.buckets[key] = b
	}

	elapsed := now.Sub(b.lastRefillTime).Seconds()
	b.tokens = math.Min(b.capacity, b.tokens+elapsed*b.refillRate)
	b.lastRefillTime = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return true, 0
	}
	need := 1 - b.tokens
	retryAfter := time.Duration((need / b.refillRate) * float64(time.Second))
	return false, retryAfter
}
