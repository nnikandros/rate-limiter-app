package ratelimiter

import (
	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Allow(string) error
	Limit() rate.Limit
}
