package goexch

import (
	"sync"
	"time"
)

// RateLimiter controls the rate of requests.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   int           // Current number of tokens
	max      int           // Maximum tokens
	interval time.Duration // Time to replenish one token
	last     time.Time     // Last time tokens were added
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(max int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:   max,
		max:      max,
		interval: interval,
		last:     time.Now(),
	}
}

// Allow checks if a request can proceed.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Replenish tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.last)
	rl.last = now

	// Add tokens for elapsed time
	rl.tokens += int(elapsed / rl.interval)
	if rl.tokens > rl.max {
		rl.tokens = rl.max
	}

	// Check if we can allow a request
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}
