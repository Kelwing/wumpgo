package rest

import (
	"golang.org/x/time/rate"
)

type MemoryConf struct {
	MaxRetries int
}

// NewMemoryRatelimiter returns a *LeakyBucketRatelimiter, but is kept for backward compatibility
func NewMemoryRatelimiter(conf *MemoryConf) *LeakyBucketRatelimiter {
	return &LeakyBucketRatelimiter{
		buckets:  make(map[string]*rate.Limiter),
		routeMap: make(map[string]string),
	}
}
