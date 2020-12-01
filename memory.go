package rest

import (
	"net/http"
	"time"
)

type MemoryConf struct {
	Authorization string
	MaxRetries    int
	UserAgent     string
}

func NewMemoryRatelimiter(conf *MemoryConf) *MemoryRatelimiter {
	return &MemoryRatelimiter{
		http:          &http.Client{Timeout: time.Second * 5},
		authorization: conf.Authorization,
		MaxRetries:    conf.MaxRetries,
		UserAgent:     conf.UserAgent,
	}
}

type MemoryRatelimiter struct {
	http          *http.Client
	authorization string
	MaxRetries    int
	UserAgent     string
}

func (r *MemoryRatelimiter) Request(method, url, contentType string, body []byte) ([]byte, error) {
	return nil, nil
}
