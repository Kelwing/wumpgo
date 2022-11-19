package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

func main() {
	logger := zerolog.New(os.Stdout)

	r, err := receiver.NewRedisReceiver(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	r.On("READY", ready)

	stop, err := r.Receive(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer stop()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func ready(c *rest.Client, r *objects.Ready) {
	fmt.Println("Ready as", r.User.Username)
}
