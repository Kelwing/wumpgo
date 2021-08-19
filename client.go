package rest

import "net/http"

type Config struct {
	Ratelimiter Ratelimiter
	Cache       Cache
}

type Client struct {
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

	resp, err := c.rateLimiter.RequestWithHeaders(req.method, req.path, req.contentType, req.body, req.headers)
	if err != nil {
		return nil, err
	}

	if req.method == "GET" && c.cache != nil {
		c.cache.Put(req.path, resp)
	}

	return resp, err
}

func (c *Client) requestNoAuth(req *request) (*DiscordResponse, error) {
	if req.headers == nil {
		req.headers = http.Header{}
	}
	if req.reason != "" && req.headers.Get(XAuditLogReasonHeader) == "" {
		req.headers.Set(XAuditLogReasonHeader, req.reason)
	}
	req.headers.Set("authorization", "")
	if req.method == "GET" && c.cache != nil {
		data, err := c.cache.Get(req.path)
		if err == nil {
			return data, nil
		}
	}
	resp, err := c.rateLimiter.RequestWithHeaders(req.method, req.path, req.contentType, req.body, req.headers)
	if err != nil {
		return nil, err
	}

	if req.method == "GET" && c.cache != nil {
		c.cache.Put(req.path, resp)
	}

	return resp, err
}
