package router

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/interactions"
	"wumpgo.dev/wumpgo/rest"
)

type RouterOption func(r *Router)

// WithInteractionsAppCmd registers this router with the given *interactions.App
func WithInteractionsApp(app *interactions.App) RouterOption {
	return func(r *Router) {
		app.CommandHandler(r.routeCommand)
		app.ComponentHandler(r.routeComponent)
		app.AutocompleteHandler(r.routeAutocomplete)
		app.ModalHandler(r.routeModal)
	}
}

// WithInitialCommands repopulates the router with the given command definitions
// WARNING: panics if any command definition is bad
func WithInitialCommands(cmds ...any) RouterOption {
	return func(r *Router) {
		for _, c := range cmds {
			r.MustRegisterCommand(c)
		}
	}
}

// WithClientCmd sets a *rest.Client on the router which will be attached
// to each CommandContext
func WithClient(c *rest.Client) RouterOption {
	return func(r *Router) {
		r.client = c
	}
}

// WithGatewayReceiverCmd configures the router to listen for interactions from the gateway
// as opposed to a webhook
func WithGatewayReceiver(rec receiver.Receiver) RouterOption {
	return func(r *Router) {
		rec.On("INTERACTION_CREATE", r.routeGatewayCommand)
		rec.On("INTERACTION_CREATE", r.routeGatewayComponent)
		rec.On("INTERACTION_CREATE", r.routeGatewayAutocomplete)
		rec.On("INTERACTION_CREATE", r.routeGatewayModal)
	}
}

// WithLogger configures the router to use the given logger instead of a noop logger
func WithLogger(l zerolog.Logger) RouterOption {
	return func(r *Router) {
		r.logger = l
	}
}
