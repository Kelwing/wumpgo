package receiver

import (
	"context"

	"wumpgo.dev/wumpgo/rest"
)

//go:generate go run ../../cmd/events ../../objects/gateway_events.go

type EventHandlerIface interface {
	New() interface{}
	Handle(context.Context, *rest.Client, interface{})
}

type EventHandler[T any] func(context.Context, *rest.Client, *T)

func (eh EventHandler[T]) New() interface{} {
	var obj T
	return &obj
}

func (eh EventHandler[T]) Handle(ctx context.Context, c *rest.Client, i interface{}) {
	if t, ok := i.(*T); ok {
		eh(ctx, c, t)
	}
}

func newHandler[T any](v EventHandler[T]) EventHandler[T] {
	return v
}
