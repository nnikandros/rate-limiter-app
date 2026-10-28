package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// RedisRateLimiter is a distributed, per-client rate limiter backed by Redis.
//
// It uses the GCRA ("leaky/token bucket") algorithm from redis_rate, enforced
// atomically in Redis via a Lua script. Because the state lives in Redis,
// multiple application instances share a single limit per client. It mirrors
// the token-bucket semantics of InMemRateLimiter: tokens refill at `rps`
// tokens/second up to a maximum of `burst` tokens.
type RedisRateLimiter struct {
	limiter *redis_rate.Limiter
	limit   redis_rate.Limit
}

// NewRedisRateLimiter creates a Redis-backed rate limiter that allows rps
// requests/sec per client with the given burst size.
func NewRedisRateLimiter(client redis.Cmdable, rps float64, burst int) *RedisRateLimiter {
	return &RedisRateLimiter{
		limiter: redis_rate.NewLimiter(client),
		limit: redis_rate.Limit{
			Rate:   int(rps),
			Burst:  burst,
			Period: time.Second,
		},
	}
}

// AllowContext reports whether a request from clientID is allowed right now.
// It returns nil if allowed, ErrRateLimited if the limit is exceeded, or a
// wrapped error if the Redis call fails.
func (rl *RedisRateLimiter) AllowContext(ctx context.Context, clientID string) error {
	res, err := rl.limiter.Allow(ctx, clientID, rl.limit)
	if err != nil {
		return fmt.Errorf("redis rate limiter: %w", err)
	}

	if res.Allowed > 0 {
		return nil
	}
	return ErrRateLimited
}

// Allow reports whether a request from clientID is allowed right now, using a
// background context. It satisfies the RateLimiter interface.
func (rl *RedisRateLimiter) Allow(clientID string) error {
	return rl.AllowContext(context.Background(), clientID)
}

// Limit returns the configured rate limit.
func (rl *RedisRateLimiter) Limit() rate.Limit {
	return rl.Limit()
}
