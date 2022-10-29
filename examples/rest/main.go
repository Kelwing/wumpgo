package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	token := flag.String("token", "", "Discord token")
	flag.Parse()

	client := rest.New(&rest.Config{
		Authorization: "Bot " + *token,
		Ratelimiter:   rest.NewLeakyBucketRatelimiter(),
		UserAgent:     "wumpgo/1.0",
	})

	me, err := client.GetCurrentUser(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get current user")
	}

	log.Info().Str("username", me.Username).Str("discriminator", me.Discriminator).Msg("")
}
