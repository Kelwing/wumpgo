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
		Expect  func(t *testing.T, c RESTClient)
	}{
		{
			Options: []RestOption{
				WithToken(objects.TokenTypeBot, "test"),
			},
			Expect: func(t *testing.T, c RESTClient) {
				cc, ok := c.(*Client)
				require.True(t, ok)
				require.Equal(t, "Bot test", cc.token)
			},
		},
		{
			Options: []RestOption{
				WithRateLimiter(rl),
			},
			Expect: func(t *testing.T, c RESTClient) {
				cc, ok := c.(*Client)
				require.True(t, ok)
				require.Equal(t, rl, cc.rateLimiter)
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
			Expect: func(t *testing.T, c RESTClient) {
				cc, ok := c.(*Client)
				require.True(t, ok)
				require.Equal(t, "wumpgo (https://wumpgo.dev, 0.0.1)", cc.userAgent.String())
			},
		},
	}

	for _, tc := range tests {
		c := New(tc.Options...)
		tc.Expect(t, c)
	}
}
