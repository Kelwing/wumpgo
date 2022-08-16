package router

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/Postcord/rest"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func TestComponentRouterCtx_DeferredMessageUpdate(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, false, "DeferredMessageUpdate"))
	assert.Equal(t, objects.ResponseDeferredMessageUpdate, x.respType)
}

func TestComponentRouterCtx_UpdateMessage(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, false, "UpdateMessage"))
	assert.Equal(t, objects.ResponseUpdateMessage, x.respType)
}

func TestComponentRouter_prep(t *testing.T) {
	x := &ComponentRouter{}
	assert.Equal(t, (map[string]any)(nil), x.routes)
	x.prep()
	assert.Equal(t, map[string]any{}, x.routes)
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

		errGeneric any
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

func jsonify(t *testing.T, x any) json.RawMessage {
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

		restClient            rest.RESTClient
		globalAllowedMentions *objects.AllowedMentions
		init                  func(t *testing.T, r *ComponentRouter, m **ModalRouter)

		interaction *objects.Interaction

		expectsErr string
		expects    *objects.InteractionResponse
	}{
		{
			name: "internal void function",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID: "/_postcord/void/1",
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
					CustomID: "/a",
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
					Content: "hello world",
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
				r.RegisterButton("/a/:content", func(ctx *ComponentRouterCtx) error {
					ctx.SetContent(ctx.Params["content"])
					return nil
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "hello",
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
			init: func(t *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
					Content: "hello world",
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
			init: func(t *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
					Content: "hello",
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
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
			init: func(_ *testing.T, r *ComponentRouter, _ **ModalRouter) {
				r.RegisterSelectMenu("/a", func(ctx *ComponentRouterCtx, _ []string) error {
					panic("wumpus fled the scene")
				})
			},
			expectsErr: "wumpus fled the scene",
		},
		{
			name: "prefer component over modal",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(t *testing.T, r *ComponentRouter, mr **ModalRouter) {
				r.RegisterButton("/a", func(ctx *ComponentRouterCtx) error {
					t.Helper()
					assert.Equal(t, dummyRestClient, ctx.RESTClient)
					ctx.SetContent("hello world")
					return nil
				})
				*mr = &ModalRouter{}
				(*mr).AddModal(&ModalContent{
					Path: "/a",
					Contents: func(ctx *ModalGenerationCtx) (name string, contents []ModalContentItem) {
						return "hello world", []ModalContentItem{}
					},
					Function: func(ctx *ModalRouterCtx) error {
						return nil
					},
				})
				(*mr).build(loaderPassthrough{})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "hello world",
				},
			},
		},
		{
			name: "modal proxy",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(t *testing.T, r *ComponentRouter, mr **ModalRouter) {
				r.RegisterButton("/b", func(ctx *ComponentRouterCtx) error {
					t.Helper()
					assert.Equal(t, dummyRestClient, ctx.RESTClient)
					ctx.SetContent("hello world")
					return nil
				})
				*mr = &ModalRouter{}
				(*mr).AddModal(&ModalContent{
					Path: "/a",
					Contents: func(ctx *ModalGenerationCtx) (name string, contents []ModalContentItem) {
						return "hello world", []ModalContentItem{
							{
								Label: "world",
							},
						}
					},
					Function: func(ctx *ModalRouterCtx) error {
						return nil
					},
				})
				(*mr).build(loaderPassthrough{})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:  objects.ComponentTypeInputText,
									Label: "world",
									Style: objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/a",
					Title:    "hello world",
				},
			},
		},
		{
			name: "modal proxy not found",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/a",
					ComponentType: objects.ComponentTypeButton,
				}),
			},
			init: func(t *testing.T, r *ComponentRouter, mr **ModalRouter) {
				*mr = &ModalRouter{}
				(*mr).build(loaderPassthrough{})
			},
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
			var m *ModalRouter
			r := &ComponentRouter{}
			if tt.init != nil {
				tt.init(t, r, &m)
			}
			builtFunc := r.build(m, loaderPassthrough{
				rest:                  dummyRestClient,
				errHandler:            setError,
				modalRouter:           m,
				globalAllowedMentions: tt.globalAllowedMentions,
				generateFrames:        false, // TODO: test frames!
			})
			resp := builtFunc(context.Background(), tt.interaction)

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
