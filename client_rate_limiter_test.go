package sourcify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRateLimiter(t *testing.T) {
	// Create a new rate limiter with max 2 actions per second
	rateLimiter := NewRateLimiter(2, time.Second)

	assert.NotNil(t, rateLimiter)
	assert.Equal(t, 2, rateLimiter.Max)
	assert.Equal(t, time.Second, rateLimiter.Duration)
}

func TestRateLimiter_Wait(t *testing.T) {
	// Create a new rate limiter with max 1 action per second
	rateLimiter := NewRateLimiter(1, time.Second)

	// Record the start time
	start := time.Now()

	// Perform 3 actions
	for i := 0; i < 3; i++ {
		rateLimiter.Wait()
	}

	// Record the end time
	end := time.Now()

	// The duration between start and end should be at least 2 seconds,
	// since the rate limiter allows only 1 action per second.
	assert.GreaterOrEqual(t, end.Sub(start).Seconds(), 2.0)
}

func TestRateLimiter_Wait_Burst(t *testing.T) {
	// Create a new rate limiter with max 5 actions per 100 milliseconds
	rateLimiter := NewRateLimiter(5, 100*time.Millisecond)

	// Record the start time
	start := time.Now()

	// Perform 5 actions, should be processed in a burst
	for i := 0; i < 5; i++ {
		rateLimiter.Wait()
	}

	// Record the end time
	end := time.Now()

	// The duration between start and end should be less than 100 milliseconds,
	// since all the actions are processed in a burst.
	assert.Less(t, end.Sub(start).Seconds(), 0.1)
}
