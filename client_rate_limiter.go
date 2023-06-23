package sourcify

import "time"

// RateLimiter represents a rate limiter that controls the rate of actions using the token bucket algorithm.
// It provides a mechanism to prevent an HTTP client from exceeding a certain rate of requests.
// The Max field represents the maximum number of actions that can be performed per 'Duration'.
// The Duration field represents the time duration for which 'Max' number of actions can be performed.
// These fields together determine the capacity of the token bucket and the rate at which tokens are added to the bucket.
// The bucket field is a channel that models the token bucket. A token is consumed from the bucket each time an action is taken.
// The capacity of the bucket determines the maximum burstiness of the actions, while the rate at which tokens are added
// to the bucket determines the sustainable average rate of actions.
type RateLimiter struct {
	// Max is the maximum number of actions that can be performed per 'Duration'.
	Max int
	// Duration is the time duration for which 'Max' number of actions can be performed.
	Duration time.Duration
	// bucket is a channel that models the token bucket. A token is consumed from the bucket each time an action is taken.
	bucket chan struct{}
}

// NewRateLimiter creates a new rate limiter.
// The rate limiter uses the token bucket algorithm to control the rate of actions.
// It initially creates a bucket of capacity 'Max' and then adds a token to the bucket every 'Duration'.
// It allows a maximum of 'Max' actions to be performed per 'Duration'.
// If an action is attempted when the bucket is empty, the action blocks until a token is added to the bucket.
// This blocking behaviour ensures that the rate of actions does not exceed the specified rate.
//
// Parameters:
// max - The maximum number of actions that can be performed per 'duration'. It is the capacity of the token bucket.
// duration - The time duration for which 'max' number of actions can be performed.
//
// Returns:
// A pointer to the created RateLimiter.
func NewRateLimiter(max int, duration time.Duration) *RateLimiter {
	bucket := make(chan struct{}, max)

	// Initially, the bucket is filled to its capacity.
	for i := 0; i < max; i++ {
		bucket <- struct{}{}
	}

	// A ticker is set up to add a token to the bucket every 'duration'.
	// If the bucket is full, the addition of a new token blocks until there is room in the bucket.
	// This ensures that the rate of actions doesn't exceed the specified rate.
	go func() {
		ticker := time.NewTicker(duration)
		for range ticker.C {
			bucket <- struct{}{}
		}
	}()

	return &RateLimiter{
		Max:      max,
		Duration: duration,
		bucket:   bucket,
	}
}

// Wait is used to perform an action with rate limiting.
// If the token bucket (i.e., 'bucket' field of RateLimiter) is empty, Wait blocks until a token is added to the bucket.
// If a token is available in the bucket, Wait consumes the token and returns immediately, allowing the action to be performed.
func (r *RateLimiter) Wait() {
	<-r.bucket
}
