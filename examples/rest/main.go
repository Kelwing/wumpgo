package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	// Parse the `token` flag from the command-line arguments.
	token := flag.String("token", "", "Discord token")
	flag.Parse()

	// Create a new `rest` client using the `token` value.
	client := rest.New(rest.WithToken(objects.TokenTypeBot, *token))

	// Get the current user using the `client`.
	me, err := client.GetCurrentUser(context.Background())
	if err != nil {
		// If there is an error, log a fatal message and exit the program.
		log.Fatal().Err(err).Msg("failed to get current user")
	}

	// Log some information about the current user.
	log.Info().Str("username", me.Username).Str("discriminator", me.Discriminator).Msg("")
}
