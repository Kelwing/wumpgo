package rest

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
)

type HTTPClient interface {
	Request(req *request) (*DiscordResponse, error)
}

type Config struct {
	Authorization string
	UserAgent     string
	Ratelimiter   Ratelimiter
	Cache         Cache
	Proxy         func(*http.Request) (*url.URL, error)
	Logger        *zerolog.Logger
}

type Client struct {
	httpClient  HTTPClient
	rateLimiter Ratelimiter
	cache       Cache
}

type request struct {
	ctx         context.Context
	method      string
	path        string
	contentType string
	body        []byte
	reason      string

	omitAuth bool

	headers http.Header

	out            interface{}
	expectedStatus []int
}

func NewRequest() *request {
	return &request{}
}

func (r *request) WithContext(ctx context.Context) *request {
	r.ctx = ctx
	return r
}

func (r *request) Expect(status ...int) *request {
	r.expectedStatus = status
	return r
}

func (r *request) OmitAuth() *request {
	r.omitAuth = true
	return r
}

func (r *request) Bind(out interface{}) *request {
	r.out = out
	return r
}

func (r *request) Method(method string) *request {
	r.method = method
	return r
}

func (r *request) Path(path string) *request {
	r.path = path
	return r
}

func (r *request) Body(body []byte) *request {
	r.body = body
	return r
}

func (r *request) ContentType(contentType string) *request {
	r.contentType = contentType
	return r
}

func (r *request) Reason(reason string) *request {
	r.reason = reason
	return r
}

func (r *request) SendRaw(c *Client) (*DiscordResponse, error) {
	if r.method == "GET" && c.cache != nil {
		data, err := c.cache.Get(r.path)
		if err == nil {
			return data, nil
		}
	}
	var resp *DiscordResponse
	var err error
	if c.rateLimiter != nil {
		resp, err = c.rateLimiter.Request(c.httpClient, r)
	} else {
		resp, err = c.httpClient.Request(r)
	}

	if err != nil {
		return nil, err
	}

	exptectedStatus := false
	for _, status := range r.expectedStatus {
		if err = resp.ExpectsStatus(status); err == nil {
			exptectedStatus = true
		}
	}
	if !exptectedStatus {
		return resp, err
	}

	if r.method == "GET" && c.cache != nil {
		// if this fails, there's not much of a recovery that can be done
		_ = c.cache.Put(r.path, resp)
	}

	return resp, err
}

func (r *request) Send(c *Client) error {
	r.headers = make(http.Header)
	if r.reason != "" && r.headers.Get(XAuditLogReasonHeader) == "" {
		r.headers.Set(XAuditLogReasonHeader, r.reason)
	}

	resp, err := r.SendRaw(c)
	if err != nil {
		return err
	}

	if r.out != nil {
		return resp.JSON(r.out)
	} else {
		return nil
	}
}

func New(config *Config) *Client {
	var client Doer
	if config.Proxy != nil {
		client = NewProxyClient(config.Proxy)
	} else {
		client = &http.Client{
			Timeout: time.Second * 5,
		}
	}
	return &Client{
		rateLimiter: config.Ratelimiter,
		httpClient: &DefaultHTTPClient{
			doer:          client,
			userAgent:     config.UserAgent,
			authorization: config.Authorization,
		},
	}
}
