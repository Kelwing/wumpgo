package interactions

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type InteractionOption func(*App)

func WithLogger(l zerolog.Logger) InteractionOption {
	return func(a *App) {
		a.logger = l
	}
}

func WithClient(c rest.RESTClient) InteractionOption {
	return func(a *App) {
		a.restClient = c
	}
}
