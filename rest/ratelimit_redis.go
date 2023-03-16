package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

var _ Ratelimiter = (*RedisRatelimiter)(nil)

type RedisConf struct {
	Options    *redis.Options
	MaxRetries int
}

func NewRedisRatelimiter(conf *RedisConf) (*RedisRatelimiter, error) {
	// TODO: Allow this client to be injected.
	r := redis.NewClient(conf.Options)
	if _, err := r.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &RedisRatelimiter{
		redis:      r,
		redsync:    redsync.New(goredis.NewPool(r)),
		MaxRetries: conf.MaxRetries,
	}, nil
}

type RedisRatelimiter struct {
	redis      *redis.Client
	redsync    *redsync.Redsync
	MaxRetries int
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

func (r *RedisRatelimiter) updateBucket(key string, resp *DiscordResponse) error {
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

func (r *RedisRatelimiter) requestLocked(httpClient HTTPClient, req *request, bucketID string, retries int) (*DiscordResponse, error) {
	if r.MaxRetries > 0 && r.MaxRetries < retries {
		return nil, ErrMaxRetriesExceeded
	}

	resp, err := httpClient.Request(req)
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
		return r.requestLocked(httpClient, req, bucketID, retries+1)
	case http.StatusBadGateway:
		return r.requestLocked(httpClient, req, bucketID, retries+1)
	}

	return resp, nil
}

func (r *RedisRatelimiter) Request(httpClient HTTPClient, req *request) (*DiscordResponse, error) {
	bucketID := getBucketID(req.path)
	mutex, err := r.acquireLock(bucketID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()
	return r.requestLocked(httpClient, req, bucketID, 0)
}
