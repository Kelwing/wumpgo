package router

import (
	"container/list"
	"errors"
)

// MiddlewareCtx is used to define the additional context that is shared between middleware.
type MiddlewareCtx struct {
	// Defines the command context.
	*CommandRouterCtx

	// Defines a list of middleware.
	middlewareList *list.List
}

// MiddlewareChainExhausted is called when the middleware chain has been exhausted.
var MiddlewareChainExhausted = errors.New("the middleware chain has been exhausted")

// Next is used to call the next function in the middleware chain.
func (m MiddlewareCtx) Next() error {
	f := m.middlewareList.Front()
	if f == nil {
		return MiddlewareChainExhausted
	}
	m.middlewareList.Remove(f)
	return f.Value.(MiddlewareFunc)(m)
}

// MiddlewareFunc is used to define a middleware function.
type MiddlewareFunc func(ctx MiddlewareCtx) error
