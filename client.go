package rest

import "net/http"

type HTTPClient interface {
	Request(req *request) (*DiscordResponse, error)
}

type Config struct {
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

	resp, err := c.rateLimiter.Request(c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		c.cache.Put(req.path, resp)
	}

	return resp, err
}
