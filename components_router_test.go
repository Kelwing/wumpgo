package router

import (
	"encoding/json"
	"errors"
	"github.com/Postcord/rest"
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

// NOT to be used, but just used to check it exists.
var dummyRestClient = &rest.Client{}

func jsonify(t *testing.T, x interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestComponentRouter_build(t *testing.T) {
	tests := []struct {
		name string

		restClient *rest.Client
		globalAllowedMentions *objects.AllowedMentions
		init func(t *testing.T, r *ComponentRouter)

		interaction *objects.Interaction

		expectsErr string
		expects    *objects.InteractionResponse
	}{
		{
			name: "internal void function",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/_postcord/void/1",
				}),
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseDeferredMessageUpdate,
			},
		},
		{
			name: "unset route",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
				}),
			},
			expects: nil,
		},

		{
			name: "button success",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterButton("/a", func(ctx *ComponentRouterCtx) error {
					if ctx.RESTClient != dummyRestClient {
						return errors.New("not dummy rest client")
					}
					ctx.SetContent("hello world")
					return nil
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content:         "hello world",
				},
			},
		},
		{
			name: "button custom param",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a/hello",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterButton("/a/:content", func(ctx *ComponentRouterCtx) error {
					ctx.SetContent(ctx.Params["content"])
					return nil
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content:         "hello",
				},
			},
		},
		{
			name: "button error",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterButton("/a", func(ctx *ComponentRouterCtx) error {
					return errors.New("wumpus fled the scene")
				})
			},
			expectsErr: "wumpus fled the scene",
		},
		{
			name: "button panic",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterButton("/a", func(ctx *ComponentRouterCtx) error {
					panic("wumpus fled the scene")
				})
			},
			expectsErr: "wumpus fled the scene",
		},

		{
			name: "select menu success",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeSelectMenu,
					Values:        []string{"a", "b", "c"},
				}),
			},
			init: func(t *testing.T, r *ComponentRouter) {
				r.RegisterSelectMenu("/a", func(ctx *ComponentRouterCtx, values []string) error {
					t.Helper()
					if ctx.RESTClient != dummyRestClient {
						return errors.New("not dummy rest client")
					}
					assert.Equal(t, []string{"a", "b", "c"}, values)
					ctx.SetContent("hello world")
					return nil
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content:         "hello world",
				},
			},
		},
		{
			name: "select menu custom param",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a/hello",
					ComponentType: objects.ComponentTypeSelectMenu,
					Values:        []string{"a", "b", "c"},
				}),
			},
			init: func(t *testing.T, r *ComponentRouter) {
				r.RegisterSelectMenu("/a/:content", func(ctx *ComponentRouterCtx, values []string) error {
					t.Helper()
					assert.Equal(t, []string{"a", "b", "c"}, values)
					ctx.SetContent(ctx.Params["content"])
					return nil
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content:         "hello",
				},
			},
		},
		{
			name: "select menu error",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeSelectMenu,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterSelectMenu("/a", func(ctx *ComponentRouterCtx, _ []string) error {
					return errors.New("wumpus fled the scene")
				})
			},
			expectsErr: "wumpus fled the scene",
		},
		{
			name: "select menu panic",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeSelectMenu,
				}),
			},
			init: func(_ *testing.T, r *ComponentRouter) {
				r.RegisterSelectMenu("/a", func(ctx *ComponentRouterCtx, _ []string) error {
					panic("wumpus fled the scene")
				})
			},
			expectsErr: "wumpus fled the scene",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle the error setter.
			var err error
			setError := func(e error) *objects.InteractionResponse {
				// Set the error.
				err = e

				// Return nil because we will not care about the response here.
				return nil
			}

			// Call the build function on the router.
			r := &ComponentRouter{}
			if tt.init != nil {
				tt.init(t, r)
			}
			builtFunc := r.build(loaderPassthrough{dummyRestClient, setError, tt.globalAllowedMentions})
			resp := builtFunc(tt.interaction)

			// Verify the error.
			if tt.expectsErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectsErr)
				return
			}
			assert.Equal(t, tt.expects, resp)
		})
	}
}
