package rest

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var _ Ratelimiter = (*RedisRatelimiter)(nil)

type RedisConf struct {
	Options       *redis.Options
	Authorization string
	MaxRetries    int
	UserAgent     string
}

func NewRedisRatelimiter(conf *RedisConf) (*RedisRatelimiter, error) {
	r := redis.NewClient(conf.Options)
	if _, err := r.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &RedisRatelimiter{
		redis:         r,
		http:          &http.Client{Timeout: time.Second * 5},
		authorization: conf.Authorization,
		redsync:       redsync.New(goredis.NewPool(r)),
		MaxRetries:    conf.MaxRetries,
		UserAgent:     conf.UserAgent,
	}, nil
}

type RedisRatelimiter struct {
	redis         *redis.Client
	redsync       *redsync.Redsync
	http          *http.Client
	authorization string
	MaxRetries    int
	UserAgent     string
}

func (r *RedisRatelimiter) getSleepTime(bucketID string) (time.Duration, error) {
	remainingRequests, err := r.redis.Get(context.Background(), bucketID+"remaining").Int()
	if err != nil && err != redis.Nil {
		return time.Duration(0), err
	}
	if remainingRequests < 1 && err != redis.Nil {
		remainingTime, err := r.redis.Get(context.Background(), bucketID+"time").Int64()
		if err != nil {
			return time.Duration(0), err
		}
		now := time.Now()
		reset := time.Unix(0, remainingTime)
		if now.Before(reset) {
			return reset.Sub(now), nil
		}
	}
	global, err := r.redis.Get(context.Background(), "global").Int64()
	if err != nil && err != redis.Nil {
		return time.Duration(0), err
	}
	now := time.Now()
	globalExpiry := time.Unix(0, global)
	if globalExpiry.After(now) {
		return globalExpiry.Sub(now), nil
	}
	return time.Duration(0), nil
}

func (r *RedisRatelimiter) acquireLock(bucketID string, opts ...redsync.Option) (*redsync.Mutex, error) {
	mutex := r.redsync.NewMutex(bucketID, opts...)
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	if waitTime, err := r.getSleepTime(bucketID); err != nil {
		return nil, err
	} else if waitTime > time.Duration(0) {
		time.Sleep(waitTime)
	}
	return mutex, nil
}

func (r *RedisRatelimiter) updateBucket(key string, resp *http.Response) error {
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
			if err = r.redis.Set(context.Background(), "global", resetAt.UnixNano(), time.Duration(0)).Err(); err != nil {
				return err
			}
		} else {
			if err = r.redis.Set(context.Background(), key+"time", resetAt.UnixNano(), time.Duration(0)).Err(); err != nil {
				return err
			}
		}
	} else if reset != "" {
		unix, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return err
		}

		resetTime := time.Unix(unix, 0).Add(time.Millisecond * 250)
		if err = r.redis.Set(context.Background(), key+"time", resetTime.UnixNano(), time.Duration(0)).Err(); err != nil {
			return err
		}
	}

	if remaining != "" {
		parsedRemaining, err := strconv.Atoi(remaining)
		if err != nil {
			return err
		}
		if err = r.redis.Set(context.Background(), key+"remaining", parsedRemaining, time.Duration(0)).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisRatelimiter) requestLocked(method, url, contentType string, body []byte, bucketID string, retries int, headers http.Header) (*DiscordResponse, error) {
	if r.MaxRetries > 0 && r.MaxRetries < retries {
		return nil, ErrMaxRetriesExceeded
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

	if r.UserAgent != "" {
		req.Header.Set("User-Agent", r.UserAgent)
	}

	if len(r.authorization) > 0 {
		req.Header.Set("authorization", r.authorization)
	}

	for k := range headers {
		// Overwrite previous headers
		v := headers.Get(k)
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}

	resp, err := r.http.Do(req)
	if err != nil {
		_ = r.updateBucket(bucketID, resp)
		return nil, err
	}

	if err = r.updateBucket(bucketID, resp); err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		delay, err := r.getSleepTime(bucketID)
		if err != nil {
			return nil, err
		}
		if delay > time.Duration(0) {
			time.Sleep(delay)
		}
		return r.requestLocked(method, url, contentType, body, bucketID, retries+1, headers)
	case http.StatusBadGateway:
		return r.requestLocked(method, url, contentType, body, bucketID, retries+1, headers)
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

func (r *RedisRatelimiter) Request(method, url, contentType string, body []byte) (*DiscordResponse, error) {
	bucketID := getBucketID(url)
	mutex, err := r.acquireLock(bucketID)
	if err != nil {
		return nil, err
	}
	defer mutex.Unlock()
	return r.requestLocked(method, url, contentType, body, bucketID, 0, nil)
}

func (r *RedisRatelimiter) RequestWithHeaders(method, url, contentType string, body []byte, headers http.Header) (*DiscordResponse, error) {
	bucketID := getBucketID(url)
	mutex, err := r.acquireLock(bucketID)
	if err != nil {
		return nil, err
	}
	defer mutex.Unlock()
	return r.requestLocked(method, url, contentType, body, bucketID, 0, headers)
}
