package rest

import (
	"bytes"
	"io"
	"io/ioutil"
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
		http:          &http.Client{Timeout: time.Second * 5},
		authorization: conf.Authorization,
		MaxRetries:    conf.MaxRetries,
		UserAgent:     conf.UserAgent,
		buckets:       make(map[string]*memoryBucket),
		global:        &global,
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
	http          *http.Client
	authorization string
	MaxRetries    int
	UserAgent     string
	buckets       map[string]*memoryBucket
	global        *int64
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

func (m *MemoryRatelimiter) requestLocked(method, url, contentType string, body []byte, bucket *memoryBucket, retries int, headers http.Header) (*DiscordResponse, error) {
	if m.MaxRetries > 0 && m.MaxRetries < retries {
		return nil, MaxRetriesExceeded
	}

	var reader io.Reader = nil
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if m.UserAgent != "" {
		req.Header.Set("User-Agent", m.UserAgent)
	}

	if len(m.authorization) > 0 {
		req.Header.Set("authorization", m.authorization)
	}

	for k := range headers {
		v := headers.Get(k)
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}

	if delay := m.getSleepTime(bucket); delay > time.Duration(0) {
		time.Sleep(delay)
	}

	resp, err := m.http.Do(req)
	if err != nil {
		_ = m.updateBucket(bucket, resp)
		return nil, err
	}

	if err = m.updateBucket(bucket, resp); err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		if delay := m.getSleepTime(bucket); delay > time.Duration(0) {
			time.Sleep(delay)
		}
		return m.requestLocked(method, url, contentType, body, bucket, retries+1, headers)
	case http.StatusBadGateway:
		return m.requestLocked(method, url, contentType, body, bucket, retries+1, headers)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &DiscordResponse{
		Body:   respBody,
		Status: resp.StatusCode,
	}, nil
}

func (m *MemoryRatelimiter) Request(method, url, contentType string, body []byte) (*DiscordResponse, error) {
	m.Lock()
	bucketID := getBucketID(url)
	bucket, ok := m.buckets[bucketID]
	if !ok {
		bucket = &memoryBucket{id: bucketID}
		m.buckets[bucketID] = bucket
	}
	bucket.Lock()
	m.Unlock()
	defer bucket.Unlock()
	return m.requestLocked(method, url, contentType, body, bucket, 0, nil)
}

func (m *MemoryRatelimiter) RequestWithHeaders(method, url, contentType string, body []byte, headers http.Header) (*DiscordResponse, error) {
	m.Lock()
	bucketID := getBucketID(url)
	bucket, ok := m.buckets[bucketID]
	if !ok {
		bucket = &memoryBucket{id: bucketID}
		m.buckets[bucketID] = bucket
	}
	bucket.Lock()
	m.Unlock()
	defer bucket.Unlock()
	return m.requestLocked(method, url, contentType, body, bucket, 0, headers)
}

func (m *MemoryRatelimiter) updateBucket(bucket *memoryBucket, resp *http.Response) error {
	if resp == nil {
		return nil
	}

	headers := resp.Header
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
