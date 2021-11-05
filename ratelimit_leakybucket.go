package rest

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

func NewLeakyBucketRatelimiter() *LeakyBucketRatelimiter {
	return &LeakyBucketRatelimiter{
		buckets:  make(map[string]*rate.Limiter),
		routeMap: make(map[string]string),
	}
}

type LeakyBucketRatelimiter struct {
	sync.RWMutex
	buckets  map[string]*rate.Limiter
	routeMap map[string]string
	logger   *zerolog.Logger
}

func (c *LeakyBucketRatelimiter) Request(httpClient HTTPClient, req *request) (*DiscordResponse, error) {
	url, err := url.Parse(req.path)
	if err != nil {
		return nil, err
	}
	bucket := c.GetBucket(url.Path)
	if bucket != nil {
		res := bucket.Reserve()
		if !res.OK() {
			return nil, ErrMaxRetriesExceeded
		}
		c.logger.Warn().Msgf("Ratelimited request to %s, waiting %v", url.Path, res.Delay())
		time.Sleep(res.Delay())
	}

	resp, err := httpClient.Request(req)
	if err != nil {
		return nil, err
	}

	c.updateFromResponse(resp.Header, req.path)

	return resp, nil
}

func (c *LeakyBucketRatelimiter) bucketExists(name string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.buckets[name]
	return ok
}

func (c *LeakyBucketRatelimiter) mappingExists(name string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.routeMap[name]
	return ok
}

func (c *LeakyBucketRatelimiter) addBucket(name string, r time.Duration, count int) {
	c.Lock()
	defer c.Unlock()
	c.logger.Debug().Msgf("Adding bucket %s with %d requests per %v seconds", name, count, r)
	c.buckets[name] = rate.NewLimiter(rate.Every(r), count)
	// Reserve a ticket since we JUST made a request
	c.buckets[name].Reserve()
}

func (c *LeakyBucketRatelimiter) addMapping(routeKey, bucket string) {
	c.Lock()
	defer c.Unlock()
	c.logger.Debug().Msgf("Adding mapping %s to bucket %s", routeKey, bucket)
	c.routeMap[routeKey] = bucket
}

func (c *LeakyBucketRatelimiter) GetBucket(path string) *rate.Limiter {
	c.RLock()
	defer c.RUnlock()
	routeKey := parseRoute(path)
	b, ok := c.buckets[c.routeMap[routeKey]]
	if !ok {
		return nil
	}
	return b
}

func (c *LeakyBucketRatelimiter) updateFromResponse(h http.Header, path string) {
	bucket := h.Get("X-RateLimit-Bucket")
	limitHeader := h.Get("X-RateLimit-Limit")
	resetAfter := h.Get("X-RateLimit-Reset-After")

	count, err := strconv.ParseInt(limitHeader, 10, 64)
	if err != nil {
		return
	}

	reset, err := strconv.ParseInt(resetAfter, 10, 64)
	if err != nil {
		return
	}

	if !c.bucketExists(bucket) {
		// Upon first request, limit and resetAfter should be actual values
		c.addBucket(bucket, time.Duration(reset)*time.Second, int(count))
	}

	routeKey := parseRoute(path)

	if !c.mappingExists(routeKey) {
		c.addMapping(routeKey, bucket)
	}
}

var snowRe = regexp.MustCompile(`\d{17,19}`)

func parseRoute(path string) string {
	splitPath := strings.Split(path, "/")
	includeNext := true
	routeKeyParts := []string{}
	for _, c := range splitPath[2:] {
		isSnowflake := snowRe.MatchString(c)
		if isSnowflake && includeNext {
			routeKeyParts = append(routeKeyParts, c)
			includeNext = false
		} else if !isSnowflake {
			routeKeyParts = append(routeKeyParts, c)
			if c == "channel" || c == "guild" || c == "webhooks" {
				includeNext = true
			}
		}
	}
	return strings.Join(routeKeyParts, ":")
}
