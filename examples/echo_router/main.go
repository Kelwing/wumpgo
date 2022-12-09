package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"wumpgo.dev/wumpgo/interactions"
	"wumpgo.dev/wumpgo/newrouter"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type Echo struct {
	Message string `discord:"message,description:Message to echo back"`
}

type EchoCaps struct {
	Message string `discord:"message,description:Message to echo back"`
}

type Log struct {
	Channel *objects.Channel `discord:",channelTypes:0"`
}

type MyCommand struct {
	Echo
	EchoCaps
	Log
}

func (e MyCommand) Description() string {
	return "Test base command"
}

func (e Echo) Description() string {
	return "Echos the message you type back to you"
}

func (e EchoCaps) Description() string {
	return "Echos the message you type back to you, but in all uppercase"
}

func (e Echo) Handle(r newrouter.CommandResponder, ctx *newrouter.CommandContext) {
	r.Content(e.Message)
}

func (e EchoCaps) Handle(r newrouter.CommandResponder, ctx *newrouter.CommandContext) {
	r.Content(strings.ToUpper(e.Message))
}

func Check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	pubKey := flag.String("pubkey", "", "Discord interactions public key")
	token := flag.String("token", "", "Discord token for creating application commands")
	flag.Parse()

	app, err := interactions.New(*pubKey)
	Check(err)

	r := newrouter.NewCommandRouter(newrouter.WithInteractionsApp(app))
	r.MustRegisterCommand(&MyCommand{})

	if *token != "" {
		fmt.Println("Creating commands with Discord")

		c := rest.New(
			rest.WithToken(objects.TokenTypeBot, *token),
			rest.WithRateLimiter(rest.NewLeakyBucketRatelimiter()),
		)

		ctx := context.Background()

		me, err := c.GetCurrentUser(ctx)
		Check(err)

		_, err = c.BulkOverwriteGlobalCommands(ctx, me.ID, r.Commands())
		Check(err)
	}

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", app)
}
