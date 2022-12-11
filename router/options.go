package router

import (
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/interactions"
	"wumpgo.dev/wumpgo/rest"
)

type CommandRouterOption func(r *CommandRouter)

// WithInteractionsApp registers this router with the given *interactions.App
func WithInteractionsApp(app *interactions.App) CommandRouterOption {
	return func(r *CommandRouter) {
		app.CommandHandler(r.routeCommand)
	}
}

// WithInitialCommands repopulates the router with the given command definitions
// WARNING: panics if any command definition is bad
func WithInitialCommands(cmds ...any) CommandRouterOption {
	return func(r *CommandRouter) {
		for _, c := range cmds {
			r.MustRegisterCommand(c)
		}
	}
}

// WithClient sets a *rest.Client on the router which will be attached
// to each CommandContext
func WithClient(c *rest.Client) CommandRouterOption {
	return func(r *CommandRouter) {
		r.client = c
	}
}

// WithGatewayReceiver configures the router to listen for interactions from the gateway
// ass opposed to a webhook
func WithGatewayReceiver(rec receiver.Receiver) CommandRouterOption {
	return func(r *CommandRouter) {
		rec.On("INTERACTION_CREATE", r.routeGatewayCommand)
	}
}
