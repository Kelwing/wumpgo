package rest

import (
	"net/http"
	"net/url"
	"time"
)

type ProxyClient struct {
	client *http.Client
}

func NewProxyClient(proxy func(*http.Request) (*url.URL, error)) *ProxyClient {
	return &ProxyClient{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: proxy,
			},
			Timeout: time.Second * 5,
		},
	}
}

func (c *ProxyClient) Do(r *http.Request) (*http.Response, error) {
	// Prevent proxy tunnelling (the proxy will upgrade this to https)
	r.URL.Scheme = "http"
	return c.client.Do(r)
}
