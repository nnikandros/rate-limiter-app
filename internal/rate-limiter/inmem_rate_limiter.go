package ratelimiter

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/time/rate"
)

// ErrRateLimited is returned by Allow when a client has exceeded its rate limit.
var ErrRateLimited = errors.New("rate limit exceeded")

// InMemRateLimiter manages per-client rate limiters in memory.
type InMemRateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*rate.Limiter
	rps      rate.Limit // requests per second
	burst    int        // max burst size
}

// NewInMemRateLimiter creates a rate limiter that allows rps requests/sec
// with the given burst size, tracked per client key.
func NewInMemRateLimiter(rps float64, burst int) *InMemRateLimiter {
	return &InMemRateLimiter{
		clients: make(map[string]*rate.Limiter),
		rps:     rate.Limit(rps),
		burst:   burst,
	}
}

// getLimiter returns the rate.Limiter for a given client, creating one if needed.
func (rl *InMemRateLimiter) getLimiter(clientID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.clients[clientID]
	if !exists {
		limiter = rate.NewLimiter(rl.rps, rl.burst)
		rl.clients[clientID] = limiter
	}
	return limiter
}

// Allow reports whether a request from clientID is allowed right now.
// It returns nil if the request is allowed, or ErrRateLimited otherwise.
func (rl *InMemRateLimiter) Allow(clientID string) error {
	if rl.getLimiter(clientID).Allow() {
		return nil
	}
	return ErrRateLimited
}

// Wait blocks until a request from clientID is allowed or ctx is cancelled.
func (rl *InMemRateLimiter) Wait(ctx context.Context, clientID string) error {
	return rl.getLimiter(clientID).Wait(ctx)
}

// Limit returns the configured rate limit.
func (rl *InMemRateLimiter) Limit() rate.Limit {
	return rl.rps
}
