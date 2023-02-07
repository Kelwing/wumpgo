package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/gateway/shard"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	token := flag.String("token", "", "Your Discord token")
	flag.Parse()

	client := rest.New(rest.WithToken(objects.TokenTypeBot, *token))

	logger := log.Logger.Level(zerolog.DebugLevel)

	gateway, err := client.GatewayBot(context.Background())
	if err != nil {
		panic(err.Error())
	}
	r := receiver.NewLocalReceiver(receiver.WithClient(client), receiver.WithLogger(logger))
	d := dispatcher.NewLocalDispatcher(r)
	r.On("READY", ready)
	r.On("GUILD_CREATE", guildCreate)
	r.On("TYPING_START", typing)

	s := shard.New(
		*token,
		shard.WithGatewayURL(gateway.URL),
		shard.WithIntents(objects.IntentsAllWithoutPrivileged),
		shard.WithDispatcher(d),
		shard.WithInitialPresence(objects.UpdatePresence{
			Activities: []objects.Activity{
				{
					Name: "I'm a bot!",
				},
			},
		}),
		shard.WithLogger(logger),
	)

	if err := s.Run(); err != nil {
		panic(err.Error())
	}
}

func ready(ctx context.Context, c *rest.Client, r *objects.Ready) {
	fmt.Println("Ready as", r.User.Username)
}

func guildCreate(ctx context.Context, c *rest.Client, g *objects.GuildCreate) {
	fmt.Println("Added to guild", g.Name)
}

func typing(ctx context.Context, c *rest.Client, t *objects.TypingStart) {
	fmt.Println(t.Member.User.Username, "started typing")
}
