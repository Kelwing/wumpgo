package router

import (
	"container/list"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareCtx_Next(t *testing.T) {
	tests := []struct {
		name string

		funcs []MiddlewareFunc
		callCount int
		expectedErr string
	}{
		{
			name:        "middleware exhausted",
			expectedErr: "the middleware chain has been exhausted",
		},
		{
			name: "first middleware",
			funcs: []MiddlewareFunc{
				func(ctx MiddlewareCtx) error {
					return errors.New("first function")
				},
			},
			expectedErr: "first function",
		},
		{
			name: "second middleware",
			funcs: []MiddlewareFunc{
				func(ctx MiddlewareCtx) error {
					return nil
				},
				func(ctx MiddlewareCtx) error {
					return errors.New("second function")
				},
			},
			callCount: 1,
			expectedErr: "second function",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := list.New()
			if tt.funcs != nil {
				for _, x := range tt.funcs {
					l.PushBack(x)
				}
			}
			m := MiddlewareCtx{middlewareList: l}
			for i := 0; i < tt.callCount; i++ {
				if err := m.Next(); err != nil {
					assert.NoError(t, err)
				}
			}
			err := m.Next()
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
