package router

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func TestComponentRouterCtx_DeferredMessageUpdate(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "DeferredMessageUpdate")
	assert.Equal(t, objects.ResponseDeferredMessageUpdate, x.respType)
}

func TestComponentRouterCtx_UpdateMessage(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "UpdateMessage")
	assert.Equal(t, objects.ResponseUpdateMessage, x.respType)
}

func TestComponentRouter_prep(t *testing.T) {
	x := &ComponentRouter{}
	assert.Equal(t, (map[string]interface{})(nil), x.routes)
	x.prep()
	assert.Equal(t, map[string]interface{}{}, x.routes)
}

func TestComponentRouter_RegisterSelectMenu(t *testing.T) {
	f := func(ctx *ComponentRouterCtx, values []string) error {
		return nil
	}
	x := &ComponentRouter{}
	x.RegisterSelectMenu("/", f)
	assert.Equal(t, reflect.Indirect(reflect.ValueOf(x.routes["/"])).Pointer(), reflect.ValueOf(f).Pointer())
}

func TestComponentRouter_RegisterButton(t *testing.T) {
	f := func(ctx *ComponentRouterCtx) error {
		return nil
	}
	x := &ComponentRouter{}
	x.RegisterButton("/", f)
	assert.Equal(t, reflect.Indirect(reflect.ValueOf(x.routes["/"])).Pointer(), reflect.ValueOf(f).Pointer())
}

func Test_ungenericError(t *testing.T) {
	tests := []struct {
		name string

		errGeneric interface{}
		expectErr  string
	}{
		{
			name:       "string",
			errGeneric: "abc",
			expectErr:  "abc",
		},
		{
			name:       "error",
			errGeneric: errors.New("abc"),
			expectErr:  "abc",
		},
		{
			name:       "other",
			errGeneric: 1,
			expectErr:  "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, ungenericError(tt.errGeneric), tt.expectErr)
		})
	}
}

func TestComponentRouter_build(t *testing.T) {
	// TODO
}
