package rest

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ProxyConf struct {
	Authorization string
	UserAgent     string
	Proxy         func(*http.Request) (*url.URL, error)
}

func NewProxyRateLimiter(conf *ProxyConf) *ProxyRateLimiter {
	tr := &http.Transport{
		Proxy: conf.Proxy,
	}

	return &ProxyRateLimiter{
		http: &http.Client{
			Timeout:   time.Second * 5,
			Transport: tr,
		},
		UserAgent:     conf.UserAgent,
		authorization: conf.Authorization,
	}
}

type ProxyRateLimiter struct {
	http          *http.Client
	authorization string
	UserAgent     string
}

func (p *ProxyRateLimiter) proxiedRequest(method, url, contentType string, body []byte, retries int, headers http.Header) (*DiscordResponse, error) {
	var reader io.Reader = nil
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = "http"

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if p.UserAgent != "" {
		req.Header.Set("User-Agent", p.UserAgent)
	}

	if len(p.authorization) > 0 {
		req.Header.Set("authorization", p.authorization)
	}

	for k := range headers {
		v := headers.Get(k)
		if v == "" {
			req.Header.Del(k)
		} else {
			req.Header.Set(k, v)
		}
	}

	resp, err := p.http.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("request error: " + string(respBody))
	}

	return &DiscordResponse{
		Body:       respBody,
		StatusCode: resp.StatusCode,
	}, nil
}

func (p *ProxyRateLimiter) Request(method, url, contentType string, body []byte) (*DiscordResponse, error) {
	return p.proxiedRequest(method, url, contentType, body, 0, nil)
}

func (p *ProxyRateLimiter) RequestWithHeaders(method, url, contentType string, body []byte, headers http.Header) (*DiscordResponse, error) {
	return p.proxiedRequest(method, url, contentType, body, 0, headers)
}
