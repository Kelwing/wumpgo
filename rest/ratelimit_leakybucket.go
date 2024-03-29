package rest

import (
	"errors"
	"math"
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

// NewLeakyBucketRatelimiter creates a new LeakyBucketRatelimiter based on rate.Limiter
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
}

func (c *LeakyBucketRatelimiter) Request(httpClient HTTPClient, req *request) (*DiscordResponse, error) {
	url, err := url.Parse(req.path)
	if err != nil {
		return nil, err
	}
	bucket, err := c.GetBucket(url.Path)
	if err == nil {
		zerolog.Ctx(req.ctx).
			Debug().
			Interface("bucket", bucket).
			Int("burst", bucket.Burst()).
			Float64("limit", float64(bucket.Limit())).
			Float64("tokens", bucket.Tokens()).
			Msgf("Processing ratelimit for %s", url.Path)
		res := bucket.Reserve()
		if !res.OK() {
			return nil, ErrMaxRetriesExceeded
		}
		zerolog.Ctx(req.ctx).Debug().
			Dur("delay", res.Delay()).
			Str("path", url.Path).
			Msg("Ratelimited request")
		time.Sleep(res.Delay())
	}

	resp, err := httpClient.Request(req)
	if err != nil {
		return nil, err
	}

	c.updateFromResponse(resp.Header, url.Path)

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

func (c *LeakyBucketRatelimiter) addBucket(name string, r float64, count float64) {
	c.Lock()
	defer c.Unlock()
	c.buckets[name] = rate.NewLimiter(rate.Every(time.Second*time.Duration(r)), int(count))
	// Reserve a ticket since we JUST made a request
	c.buckets[name].ReserveN(time.Now(), int(count))
}

func (c *LeakyBucketRatelimiter) addMapping(routeKey, bucket string) {
	c.Lock()
	defer c.Unlock()
	c.routeMap[routeKey] = bucket
}

func (c *LeakyBucketRatelimiter) GetBucket(path string) (*rate.Limiter, error) {
	c.RLock()
	defer c.RUnlock()

	routeKey := parseRoute(path)
	bucketName, ok := c.routeMap[routeKey]
	if !ok {
		return nil, errors.New("no bucket mapping found")
	}

	if bucket, ok := c.buckets[bucketName]; ok {
		return bucket, nil
	}

	return nil, errors.New("no bucket found")
}

func (c *LeakyBucketRatelimiter) updateFromResponse(h http.Header, path string) {
	bucket := h.Get("X-RateLimit-Bucket")
	limitHeader := h.Get("X-RateLimit-Limit")
	resetAfter := h.Get("X-RateLimit-Reset-After")

	count, err := strconv.ParseInt(limitHeader, 10, 64)
	if err != nil {
		return
	}

	reset, err := strconv.ParseFloat(resetAfter, 64)
	if err != nil {
		return
	}

	reset = math.Round(reset/0.25) * 0.25

	if !c.bucketExists(bucket) {
		// Upon first request, limit and resetAfter should be actual values
		c.addBucket(bucket, reset, float64(count))
	}

	routeKey := parseRoute(path)

	if !c.mappingExists(routeKey) {
		c.addMapping(routeKey, bucket)
	}
}

var snowRe = regexp.MustCompile(`\d{17,19}`)

func parseRoute(path string) string {
	path, _, _ = strings.Cut(path, "?")
	splitPath := strings.Split(path, "/")
	includeNext := true
	routeKeyParts := []string{}
	for _, c := range splitPath[3:] {
		isSnowflake := snowRe.MatchString(c)
		if isSnowflake && includeNext {
			routeKeyParts = append(routeKeyParts, c)
			includeNext = false
		} else if !isSnowflake {
			routeKeyParts = append(routeKeyParts, c)
			if c == "channels" || c == "guilds" || c == "webhooks" {
				includeNext = true
			}
		}
	}
	return strings.Join(routeKeyParts, ":")
}
