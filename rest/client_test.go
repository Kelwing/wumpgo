package rest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"wumpgo.dev/wumpgo/objects"
)

func TestNew(t *testing.T) {
	rl := NewLeakyBucketRatelimiter()
	tests := []struct {
		Options []RestOption
		Expect  func(t *testing.T, c *Client)
	}{
		{
			Options: []RestOption{
				WithToken(objects.TokenTypeBot, "test"),
			},
			Expect: func(t *testing.T, c *Client) {
				require.Equal(t, "Bot test", c.token)
			},
		},
		{
			Options: []RestOption{
				WithRateLimiter(rl),
			},
			Expect: func(t *testing.T, c *Client) {
				require.Equal(t, rl, c.rateLimiter)
			},
		},
		{
			Options: []RestOption{
				WithUserAgent(&UserAgent{
					Name:    "wumpgo",
					URL:     "https://wumpgo.dev",
					Version: "0.0.1",
				}),
			},
			Expect: func(t *testing.T, c *Client) {
				require.Equal(t, "wumpgo (https://wumpgo.dev, 0.0.1)", c.userAgent.String())
			},
		},
	}

	for _, tc := range tests {
		c := New(tc.Options...)
		tc.Expect(t, c)
	}
}
