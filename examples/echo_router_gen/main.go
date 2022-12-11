package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"wumpgo.dev/wumpgo/interactions"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
	"wumpgo.dev/wumpgo/router"
)

// go:generate cmdgen

// Echo godoc
// @Description Echos the message you type back to you
type Echo struct {
	Message string `discord:"message,description:Message to echo back"`
}

// EchoCaps godoc
// @Description Echos the message you type back to you, but in all uppercase
// @Option.Message.Name.es_MX mensaje
// @Option.Message.Name.en_US message
// @Permissions ManageRoles
type EchoCaps struct {
	Message string `discord:"message,description:Message to echo back"`
}

// Log godoc
// @Description Tests a channel argument
type Log struct {
	Channel *objects.Channel `discord:",channelTypes:0"`
}

// MyCommand godoc
// @Name testcommand
// @Description Test base command
// @Name.en_US testcommand
// @Name.es_MX commandodepreueba
// @Type ChatInput
type MyCommand struct {
	Echo
	EchoCaps
	Log
}

func (e Echo) Handle(r router.CommandResponder, ctx *router.CommandContext) {
	r.Content(e.Message)
}

func (e EchoCaps) Handle(r router.CommandResponder, ctx *router.CommandContext) {
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

	r := router.NewCommandRouter(router.WithInteractionsApp(app))
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
