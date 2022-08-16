package rest

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetSleepTime(t *testing.T) {
	rateLimiter := NewMemoryRatelimiter(&MemoryConf{MaxRetries: 3})
	header := http.Header{}

	waitSeconds := 15

	header.Add("X-RateLimit-Limit", "3")
	header.Add("X-RateLimit-Remaining", "0")
	header.Add("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second*time.Duration(waitSeconds)).Unix()))
	header.Add("X-RateLimit-Reset-After", fmt.Sprintf("%d", waitSeconds))
	header.Add("X-RateLimit-Bucket", "e06f83c33559dfd4dc34f5666fdfa1d3")

	resp := DiscordResponse{
		StatusCode: 200,
		Header:     header,
		Body:       []byte("{}"),
	}

	bucket := memoryBucket{}
	if err := rateLimiter.updateBucket(&bucket, &resp); err != nil {
		t.Fail()
	}
	expectedDelay := time.Second * time.Duration(waitSeconds)
	if delay := rateLimiter.getSleepTime(&bucket); delay < (expectedDelay-time.Second) || delay > (expectedDelay+time.Second) {
		t.Errorf("Expected delay of %v, got %v", expectedDelay, delay)
	}
}
