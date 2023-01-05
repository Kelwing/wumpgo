package loader

import (
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/router"
)

type LoaderOption func(*Loader)

func WithRouter(r *router.Router) LoaderOption {
	return func(l *Loader) {
		l.router = r
	}
}

func WithGatewayReceiver(r receiver.Receiver) LoaderOption {
	return func(l *Loader) {
		l.rec = r
	}
}
