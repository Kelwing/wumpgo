package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/objects"
)

type RestOption func(*Client)

func WithToken(t objects.TokenType, token string) RestOption {
	return func(c *Client) {
		c.token = fmt.Sprintf("%s %s", t, token)
	}
}

func WithRateLimiter(r Ratelimiter) RestOption {
	return func(c *Client) {
		c.rateLimiter = r
	}
}

func WithCache(cache Cache) RestOption {
	return func(c *Client) {
		c.cache = cache
	}
}

func WithLogger(l zerolog.Logger) RestOption {
	return func(c *Client) {
		c.logger = l
	}
}

func WithProxy(p func(*http.Request) (*url.URL, error)) RestOption {
	return func(c *Client) {
		c.proxy = p
	}
}

func WithUserAgent(a *UserAgent) RestOption {
	return func(c *Client) {
		c.userAgent = a
	}
}
