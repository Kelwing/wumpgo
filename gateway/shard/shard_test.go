package shard

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func newShard() (*Shard, error) {
	logger := log.Level(zerolog.DebugLevel)
	token := os.Getenv("TOKEN")
	// Create a REST client with ratelimiting
	client := rest.New(
		rest.WithToken(objects.TokenTypeBot, token),
		rest.WithRateLimiter(rest.NewLeakyBucketRatelimiter()),
		rest.WithUserAgent(&rest.UserAgent{
			Name:    "wumpgo-test",
			URL:     "wumpgo.dev",
			Version: "0.0.1",
		}),
		rest.WithLogger(logger),
	)

	gateway, err := client.GatewayBot(context.Background())
	if err != nil {
		return nil, err
	}

	d := dispatcher.NewNOOPDispatcher(dispatcher.WithLogger(&logger))

	s := New(
		token,
		WithGatewayURL(gateway.URL),
		WithIntents(objects.IntentsNone),
		WithDispatcher(d),
		WithLogger(logger),
	)

	return s, nil
}

func TestClose(t *testing.T) {
	s, err := newShard()
	require.NoError(t, err)

	shardErr := make(chan error)

	go func(s *Shard) {
		err := s.Run()
		shardErr <- err
	}(s)

	time.Sleep(time.Second * 5)
	s.Close()
	err = <-shardErr
	var gwErr ShardError
	require.ErrorAs(t, err, &gwErr)
}
