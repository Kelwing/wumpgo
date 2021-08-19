package rest

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type MemoryConf struct {
	Authorization string
	MaxRetries    int
	UserAgent     string
}

func NewMemoryRatelimiter(conf *MemoryConf) *MemoryRatelimiter {
	global := int64(0)
	return &MemoryRatelimiter{
		MaxRetries: conf.MaxRetries,
		buckets:    make(map[string]*memoryBucket),
		global:     &global,
	}
}

var _ Ratelimiter = (*MemoryRatelimiter)(nil)

type memoryBucket struct {
	sync.Mutex
	resetTime time.Time
	remaining int
	id        string
}

type MemoryRatelimiter struct {
	sync.Mutex
	MaxRetries int
	buckets    map[string]*memoryBucket
	global     *int64
}

func (m *MemoryRatelimiter) getSleepTime(bucket *memoryBucket) time.Duration {
	now := time.Now()
	if bucket.remaining < 1 && now.Before(bucket.resetTime) {
		return now.Sub(bucket.resetTime)
	}
	globalReset := atomic.LoadInt64(m.global)
	globalTime := time.Unix(0, globalReset)
	if now.Before(globalTime) {
		return now.Sub(globalTime)
	}
	return time.Duration(0)
}

func (m *MemoryRatelimiter) requestLocked(httpClient HTTPClient, r *request, bucket *memoryBucket, retries int) (*DiscordResponse, error) {
	if m.MaxRetries > 0 && m.MaxRetries < retries {
		return nil, ErrMaxRetriesExceeded
	}

	if delay := m.getSleepTime(bucket); delay > time.Duration(0) {
		time.Sleep(delay)
	}

	resp, err := httpClient.Request(r)
	if err != nil {
		_ = m.updateBucket(bucket, resp.Header)
		return nil, err
	}

	if err = m.updateBucket(bucket, resp.Header); err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		if delay := m.getSleepTime(bucket); delay > time.Duration(0) {
			time.Sleep(delay)
		}
		return m.requestLocked(httpClient, r, bucket, retries+1)
	case http.StatusBadGateway:
		return m.requestLocked(httpClient, r, bucket, retries+1)
	}

	return resp, nil
}

func (m *MemoryRatelimiter) Request(httpClient HTTPClient, req *request) (*DiscordResponse, error) {
	m.Lock()
	bucketID := getBucketID(req.path)
	bucket, ok := m.buckets[bucketID]
	if !ok {
		bucket = &memoryBucket{id: bucketID}
		m.buckets[bucketID] = bucket
	}
	bucket.Lock()
	m.Unlock()
	defer bucket.Unlock()
	return m.requestLocked(httpClient, req, bucket, 0)
}

func (m *MemoryRatelimiter) updateBucket(bucket *memoryBucket, headers http.Header) error {
	remaining := headers.Get("X-RateLimit-Remaining")
	reset := headers.Get("X-RateLimit-Reset")
	global := headers.Get("X-RateLimit-Global")
	resetAfter := headers.Get("X-RateLimit-Reset-After")

	if resetAfter != "" {
		parsedAfter, err := strconv.ParseFloat(resetAfter, 64)
		if err != nil {
			return err
		}

		resetAt := time.Now().Add(time.Duration(parsedAfter) * time.Millisecond)

		if global != "" {
			atomic.StoreInt64(m.global, resetAt.UnixNano())
		} else {
			bucket.resetTime = resetAt
		}
	} else if reset != "" {
		unix, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return err
		}

		resetTime := time.Unix(unix, 0).Add(time.Millisecond * 250)
		bucket.resetTime = resetTime
	}

	if remaining != "" {
		parsedRemaining, err := strconv.Atoi(remaining)
		if err != nil {
			return err
		}
		bucket.remaining = parsedRemaining
	}
	return nil
}
