package main

import (
	"context"
	"flag"

	"github.com/go-redis/redis/v8"
	"wumpgo.dev/wumpgo/gateway"
	"wumpgo.dev/wumpgo/gateway/shard"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	token := flag.String("token", "", "Your Discord token")
	flag.Parse()

	client := rest.New(rest.WithToken(objects.TokenTypeBot, *token))

	gw, err := client.GatewayBot(context.Background())
	if err != nil {
		panic(err.Error())
	}
	d, err := gateway.NewRedisDispatcher(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err != nil {
		panic(err.Error())
	}

	s := gateway.NewShard(
		*token,
		shard.WithGatewayURL(gw.URL),
		shard.WithIntents(objects.IntentsAllWithoutPrivileged|objects.IntentsGuildMessages),
		shard.WithDispatcher(d),
		shard.WithInitialPresence(objects.UpdatePresence{
			Activities: []objects.Activity{
				{
					Name: "I'm a bot!",
				},
			},
		}),
	)

	if err := s.Run(); err != nil {
		panic(err.Error())
	}
}
