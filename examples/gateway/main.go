package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/gateway/shard"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	token := flag.String("token", "", "Your Discord token")
	flag.Parse()

	client := rest.New(&rest.Config{
		Authorization: "Bot " + *token,
		Ratelimiter:   rest.NewLeakyBucketRatelimiter(),
	})

	gateway, err := client.GatewayBot(context.Background())
	if err != nil {
		panic(err.Error())
	}

	d := dispatcher.NewLocalDispatcher(client, zerolog.Nop())
	d.On("READY", ready)
	d.On("GUILD_CREATE", guildCreate)

	s := shard.New(
		*token,
		shard.WithGatewayURL(gateway.URL),
		shard.WithIntents(objects.IntentsAllWithoutPrivileged),
		shard.WithDispatcher(d),
	)

	if err := s.Run(); err != nil {
		panic(err.Error())
	}
}

func ready(c *rest.Client, r *objects.Ready) {
	fmt.Println("Ready as", r.User.Username)
}

func guildCreate(c *rest.Client, g *objects.GuildCreate) {
	fmt.Println("Added to guild", g.Name)
}
