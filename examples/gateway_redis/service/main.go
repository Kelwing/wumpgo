package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	logger := zerolog.New(os.Stdout)

	r, err := gateway.NewRedisReceiver(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	r.On(ready)
	r.On(guildCreate)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := r.Run(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to run")
	}
}

func ready(ctx context.Context, c *rest.Client, r *objects.Ready) {
	fmt.Println("Ready as", r.User.Username)
}

func guildCreate(ctx context.Context, c *rest.Client, g *objects.GuildCreate) {
	fmt.Println("Added to guild", g.Name)
}
