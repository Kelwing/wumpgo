package rest

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Request(req *request) (*DiscordResponse, error)
}

type Config struct {
	Token       string
	UserAgent   string
	Ratelimiter Ratelimiter
	Cache       Cache
}

type Client struct {
	httpClient  HTTPClient
	rateLimiter Ratelimiter
	cache       Cache
}

func New(config *Config) *Client {
	return &Client{
		rateLimiter: config.Ratelimiter,
		httpClient: &DefaultHTTPClient{
			doer: &http.Client{
				Timeout: time.Second * 5,
			},
			userAgent:     config.UserAgent,
			authorization: config.Token,
		},
	}
}

func (c *Client) request(req *request) (*DiscordResponse, error) {
	req.headers = make(http.Header)
	if req.reason != "" && req.headers.Get(XAuditLogReasonHeader) == "" {
		req.headers.Set(XAuditLogReasonHeader, req.reason)
	}

	if req.method == "GET" && c.cache != nil {
		data, err := c.cache.Get(req.path)
		if err == nil {
			return data, nil
		}
	}
  
	var resp *DiscordResponse
	var err error
	if c.rateLimiter != nil {
		resp, err = c.rateLimiter.Request(c.httpClient, req)
	} else {
		resp, err = c.httpClient.Request(req)
  }

	if req.method == "GET" && c.cache != nil {
		c.cache.Put(req.path, resp)
	}

	return resp, err
}
