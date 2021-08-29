package router

import (
	"reflect"
	"testing"

	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"github.com/stretchr/testify/assert"
)

func TestRouterBuilder(t *testing.T) {
	assert.Equal(t, &loaderBuilder{}, RouterLoader())
}

func TestLoaderBuilder_AllowedMentions(t *testing.T) {
	tests := []struct{
		name string

		allowedMentions *objects.AllowedMentions
	}{
		{
			name: "value",
			allowedMentions: &objects.AllowedMentions{
				Parse: []string{"a", "b", "c"},
			},
		},
		{
			name: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := RouterLoader().(*loaderBuilder)
			if tt.allowedMentions == nil {
				// Ensure it is nil-ed.
				l.globalAllowedMentions	= &objects.AllowedMentions{}
			}
			l.AllowedMentions(tt.allowedMentions)
			assert.Equal(t, tt.allowedMentions, l.globalAllowedMentions)
		})
	}
}

func TestLoaderBuilder_ErrorHandler(t *testing.T) {
	tests := []struct{
		name string

		handler func(error) *objects.InteractionResponse
	}{
		{
			name: "value",
			handler: func(error) *objects.InteractionResponse {
				return nil
			},
		},
		{
			name: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := RouterLoader().(*loaderBuilder)
			if tt.handler == nil {
				// Ensure it is nil-ed.
				l.errHandler = func(err error) *objects.InteractionResponse {
					return nil
				}
			}
			l.ErrorHandler(tt.handler)
			assert.Equal(t, reflect.ValueOf(tt.handler).Pointer(), reflect.ValueOf(l.errHandler).Pointer())
		})
	}
}

func TestLoaderBuilder_ComponentRouter(t *testing.T) {
	tests := []struct{
		name string

		value *ComponentRouter
	}{
		{
			name: "value",
			value: &ComponentRouter{},
		},
		{
			name: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := RouterLoader().(*loaderBuilder)
			if tt.value == nil {
				// Ensure it is nil-ed.
				l.components = &ComponentRouter{}
			}
			l.ComponentRouter(tt.value)
			assert.Equal(t, reflect.ValueOf(tt.value).Pointer(), reflect.ValueOf(l.components).Pointer())
		})
	}
}

func TestLoaderBuilder_CommandRouter(t *testing.T) {
	tests := []struct{
		name string

		value *CommandRouter
	}{
		{
			name: "value",
			value: &CommandRouter{},
		},
		{
			name: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := RouterLoader().(*loaderBuilder)
			if tt.value == nil {
				// Ensure it is nil-ed.
				l.commands = &CommandRouter{}
			}
			l.CommandRouter(tt.value)
			assert.Equal(t, reflect.ValueOf(tt.value).Pointer(), reflect.ValueOf(l.commands).Pointer())
		})
	}
}

func TestLoaderBuilder_CombinedRouter(t *testing.T) {
	tests := []struct{
		name string

		value *CombinedRouter
	}{
		{
			name: "value",
			value: &CombinedRouter{},
		},
		{
			name: "nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := RouterLoader().(*loaderBuilder)
			if tt.value == nil {
				// Ensure it is nil-ed.
				l.commands = &CommandRouter{}
				l.components = &ComponentRouter{}
			}
			l.CombinedRouter(tt.value)
			if tt.value == nil {
				assert.Equal(t, (*ComponentRouter)(nil), l.components)
				assert.Equal(t, (*CommandRouter)(nil), l.commands)
			} else {
				assert.Equal(t, reflect.ValueOf(&tt.value.ComponentRouter).Pointer(), reflect.ValueOf(l.components).Pointer())
				assert.Equal(t, reflect.ValueOf(&tt.value.CommandRouter).Pointer(), reflect.ValueOf(l.commands).Pointer())
			}
		})
	}
}

type fakeBuildHandlerAccepter struct {
	componentHandler, commandHandler interactions.HandlerFunc
}

func (f *fakeBuildHandlerAccepter) ComponentHandler(handler interactions.HandlerFunc) {
	f.componentHandler = handler
}

func (f *fakeBuildHandlerAccepter) CommandHandler(handler interactions.HandlerFunc) {
	f.commandHandler = handler
}

func (fakeBuildHandlerAccepter) Rest() *rest.Client {
	return nil
}

func TestLoaderBuilder_Build(t *testing.T) {
	tests := []struct {
		name string

		allowedMentions *objects.AllowedMentions
		components *ComponentRouter
		commands *CommandRouter
		errHandler func(error) *objects.InteractionResponse
	}{
		{
			name: "all nil",
		},
		{
			name: "command router nil",
			components: &ComponentRouter{},
		},
		{
			name: "component router nil",
			commands: &CommandRouter{},
		},
		{
			name: "error handler present",
			components: &ComponentRouter{},
			commands: &CommandRouter{},
			errHandler: func(err error) *objects.InteractionResponse {
				return nil
			},
		},
		{
			name: "allowed mentions present",
			components: &ComponentRouter{},
			commands: &CommandRouter{},
			allowedMentions: &objects.AllowedMentions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &fakeBuildHandlerAccepter{}
			RouterLoader().
				AllowedMentions(tt.allowedMentions).
				ComponentRouter(tt.components).
				CommandRouter(tt.commands).
				ErrorHandler(tt.errHandler).
				Build(app)
			if tt.components == nil {
				assert.Nil(t, app.componentHandler)
			} else {
				assert.NotNil(t, app.componentHandler)
			}
			if tt.commands == nil {
				assert.Nil(t, app.commandHandler)
			} else {
				assert.NotNil(t, app.commandHandler)
			}
		})
	}
}
