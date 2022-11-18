package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	token := flag.String("token", "", "Discord token")
	flag.Parse()

	client := rest.New(rest.WithToken(objects.TokenTypeBot, *token))

	me, err := client.GetCurrentUser(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get current user")
	}

	log.Info().Str("username", me.Username).Str("discriminator", me.Discriminator).Msg("")
}
