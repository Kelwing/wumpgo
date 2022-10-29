package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/interactions"
	"wumpgo.dev/wumpgo/objects"
)

func main() {
	key := flag.String("public_key", "", "Discord public key")
	flag.Parse()

	app, err := interactions.New(&interactions.Config{
		PublicKey: *key,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create interactions client")
	}

	app.CommandHandler(myCommandHandler)
	app.ComponentHandler(myComponentHandler)
}

func myCommandHandler(ctx context.Context, i *objects.Interaction) *objects.InteractionResponse {
	return &objects.InteractionResponse{
		Type: objects.ResponseChannelMessageWithSource,
		Data: &objects.InteractionApplicationCommandCallbackData{
			Content: "Hello world!",
			Flags:   objects.MsgFlagEphemeral,
			Components: []*objects.Component{
				{
					Type: objects.ComponentTypeActionRow,
					Components: []*objects.Component{
						{
							Type:     objects.ComponentTypeButton,
							CustomID: "btn_hello",
							Label:    "Hello!",
							Style:    objects.ButtonStyleSuccess,
						},
					},
				},
			},
		},
	}
}

func myComponentHandler(ctx context.Context, i *objects.Interaction) *objects.InteractionResponse {
	return &objects.InteractionResponse{
		Type: objects.ResponseUpdateMessage,
		Data: &objects.InteractionApplicationCommandCallbackData{
			Content: "Goodbye world!",
		},
	}
}
