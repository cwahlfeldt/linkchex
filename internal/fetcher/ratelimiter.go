package fetcher

import (
	"sync"
	"time"
)

// RateLimiter controls the rate of requests
type RateLimiter struct {
	requestsPerSecond float64
	ticker            *time.Ticker
	tokens            chan struct{}
	stop              chan struct{}
	once              sync.Once
}

// NewRateLimiter creates a new rate limiter
// requestsPerSecond: number of requests allowed per second (0 = unlimited)
func NewRateLimiter(requestsPerSecond float64) *RateLimiter {
	if requestsPerSecond <= 0 {
		// No rate limiting
		return &RateLimiter{
			requestsPerSecond: 0,
		}
	}

	// Calculate interval between requests
	interval := time.Duration(float64(time.Second) / requestsPerSecond)

	rl := &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		ticker:            time.NewTicker(interval),
		tokens:            make(chan struct{}, 1),
		stop:              make(chan struct{}),
	}

	// Start the token generator
	go rl.generate()

	// Add initial token
	rl.tokens <- struct{}{}

	return rl
}

// generate continuously generates tokens at the specified rate
func (rl *RateLimiter) generate() {
	for {
		select {
		case <-rl.stop:
			return
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Token channel is full, skip
			}
		}
	}
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	if rl.requestsPerSecond == 0 {
		// No rate limiting
		return
	}
	<-rl.tokens
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	if rl.requestsPerSecond == 0 {
		return
	}
	rl.once.Do(func() {
		close(rl.stop)
		rl.ticker.Stop()
	})
}
