package rest

import "net/http"

type Config struct {
	Ratelimiter Ratelimiter
}

type Client struct {
	rateLimiter Ratelimiter
}

func New(config *Config) *Client {
	return &Client{
		rateLimiter: config.Ratelimiter,
	}
}

func (c *Client) request(req *request) (*DiscordResponse, error) {
	if req.reason != "" && req.headers.Get(XAuditLogReasonHeader) == "" {
		req.headers.Set(XAuditLogReasonHeader, req.reason)
	}
	return c.rateLimiter.RequestWithHeaders(req.method, req.path, req.contentType, req.body, req.headers)
}

func (c *Client) requestNoAuth(req *request) (*DiscordResponse, error) {
	if req.headers == nil {
		req.headers = http.Header{}
	}
	if req.reason != "" && req.headers.Get(XAuditLogReasonHeader) == "" {
		req.headers.Set(XAuditLogReasonHeader, req.reason)
	}
	req.headers.Set("authorization", "")
	return c.rateLimiter.RequestWithHeaders(req.method, req.path, req.contentType, req.body, req.headers)
}
