package receiver

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type ReceiverOption func(*eventRouter)

func WithClient(c *rest.Client) ReceiverOption {
	return func(e *eventRouter) {
		e.client = c
	}
}

func WithLogger(l zerolog.Logger) ReceiverOption {
	return func(e *eventRouter) {
		e.log = l
	}
}

func WithErrorHandler(h func(error)) ReceiverOption {
	return func(e *eventRouter) {
		e.errHandler = h
	}
}

func WithGroupName(name string) ReceiverOption {
	return func(e *eventRouter) {
		e.groupName = name
	}
}
