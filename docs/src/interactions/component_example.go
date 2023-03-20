package example

import (
	"time"

	"wumpgo.dev/wumpgo/router"
)

//go:generate wumpgoctl gen

// Ping is a simple command that returns the string "Pong!"
// @Name ping
// @Description Check to see if the bot is responding
// @DM true
type Ping struct{}

func (p Ping) Handle(r router.CommandResponder, ctx *router.CommandContext) {
	pingTime := time.Since(ctx.Interaction().ID.CreatedAt().Time)

	r.Contentf("Pong! %s", pingTime.String()).
		View(router.NewView().Add(
			router.NewButton("/ping").Label("Ping Again"),
		))
}

func PingAgain(r router.ComponentResponder, ctx *router.ComponentContext) {
	pingTime := time.Since(ctx.Interaction().ID.CreatedAt().Time)

	r.UpdateMessage().Contentf("Pong! %s", pingTime.String()).
		View(router.NewView().Add(
			router.NewButton("/ping").Label("Ping Again"),
		))
}

func Register(r *router.Router) {
	r.AddHandler("/ping", PingAgain)
}
