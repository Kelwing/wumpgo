package router

import (
	"context"
	"errors"
	"testing"

	"wumpgo.dev/wumpgo/rest"

	"wumpgo.dev/wumpgo/objects"
	"github.com/stretchr/testify/assert"
)

func TestModalRouter_prep(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		m := ModalRouter{}
		m.prep()
		assert.NotNil(t, m.routes)
	})
	t.Run("map exists", func(t *testing.T) {
		m := ModalRouter{routes: map[string]*ModalContent{"a": nil}}
		m.prep()
		assert.Equal(t, map[string]*ModalContent{"a": nil}, m.routes)
	})
}

func Test_uint2IntPtr(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		assert.Nil(t, uint2IntPtr(0))
	})
	t.Run("non-zero", func(t *testing.T) {
		x := 1
		assert.Equal(t, &x, uint2IntPtr(1))
	})
}

type builder interface {
	ResponseData() *objects.InteractionApplicationCommandCallbackData
	buildResponse(component bool, errorHandler ErrorHandler, globalAllowedMentions *objects.AllowedMentions) *objects.InteractionResponse
}

type fakeCtx struct{}

func (fakeCtx) ResponseData() *objects.InteractionApplicationCommandCallbackData {
	panic("should not be called")
}

func (fakeCtx) buildResponse(_ bool, _ ErrorHandler, _ *objects.AllowedMentions) *objects.InteractionResponse {
	panic("should not be called")
}

func intPtr(x int) *int {
	return &x
}

func TestModalRouter_SendModalResponse(t *testing.T) {
	tests := []struct {
		name string

		ctx      builder
		init     func(*ModalRouter)
		path     string
		wantsErr string
		wants    *objects.InteractionResponse
	}{
		{
			name:     "no modal path",
			wantsErr: "modal path not found",
		},
		{
			name: "modal ctx",
			path: "/a",
			ctx:  &ModalRouterCtx{},
			init: func(r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("should not be called")
					},
					Function: func(_ *ModalRouterCtx) error {
						panic("should not be called")
					},
				})
			},
			wantsErr: "multiple modal responses",
		},
		{
			name: "command ctx",
			path: "/a",
			ctx:  &CommandRouterCtx{},
			init: func(r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						return "Hello World", []ModalContentItem{
							{
								Short:       true,
								Label:       "abc",
								Key:         "def",
								Placeholder: "ghi",
								Value:       "jkl",
								Required:    true,
								MinLength:   1,
								MaxLength:   4,
							},
							{
								Short:       false,
								Label:       "abclol",
								Key:         "def",
								Placeholder: "ghi",
								Value:       "jkl",
								Required:    true,
								MinLength:   0,
								MaxLength:   0,
							},
						}
					},
					Function: func(_ *ModalRouterCtx) error {
						panic("should not be called")
					},
				})
			},
			wants: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:        objects.ComponentTypeInputText,
									Label:       "abc",
									CustomID:    "def",
									Placeholder: "ghi",
									Value:       "jkl",
									MinLength:   intPtr(1),
									MaxLength:   intPtr(4),
									Required:    true,
									Style:       objects.ButtonStyle(objects.TextStyleShort),
								},
							},
						},
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:        objects.ComponentTypeInputText,
									Label:       "abclol",
									CustomID:    "def",
									Placeholder: "ghi",
									Value:       "jkl",
									Required:    true,
									Style:       objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/a",
					Title:    "Hello World",
				},
			},
		},
		{
			name: "component ctx",
			path: "/a",
			ctx:  &ComponentRouterCtx{},
			init: func(r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						return "Hello World", []ModalContentItem{
							{
								Short:       true,
								Label:       "abc",
								Key:         "def",
								Placeholder: "ghi",
								Value:       "jkl",
								Required:    true,
								MinLength:   1,
								MaxLength:   4,
							},
							{
								Short:       false,
								Label:       "abclol",
								Key:         "def",
								Placeholder: "ghi",
								Value:       "jkl",
								Required:    true,
								MinLength:   0,
								MaxLength:   0,
							},
						}
					},
					Function: func(_ *ModalRouterCtx) error {
						panic("should not be called")
					},
				})
			},
			wants: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:        objects.ComponentTypeInputText,
									Label:       "abc",
									CustomID:    "def",
									Placeholder: "ghi",
									Value:       "jkl",
									MinLength:   intPtr(1),
									MaxLength:   intPtr(4),
									Required:    true,
									Style:       objects.ButtonStyle(objects.TextStyleShort),
								},
							},
						},
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:        objects.ComponentTypeInputText,
									Label:       "abclol",
									CustomID:    "def",
									Placeholder: "ghi",
									Value:       "jkl",
									Required:    true,
									Style:       objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/a",
					Title:    "Hello World",
				},
			},
		},
		{
			name: "unknown ctx",
			path: "/a",
			ctx:  &fakeCtx{},
			init: func(r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("should not be called")
					},
					Function: func(_ *ModalRouterCtx) error {
						panic("should not be called")
					},
				})
			},
			wantsErr: "unknown context type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ModalRouter{}
			if tt.init != nil {
				tt.init(r)
			}
			r.build(loaderPassthrough{})
			err := r.SendModalResponse(tt.ctx, tt.path)
			if err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wants, tt.ctx.buildResponse(false, nil, nil))
			} else {
				assert.EqualError(t, err, tt.wantsErr)
			}
		})
	}
}

func TestModalRouter_build(t *testing.T) {
	tests := []struct {
		name string

		restClient            rest.RESTClient
		globalAllowedMentions *objects.AllowedMentions
		init                  func(t *testing.T, m *ModalRouter)

		interaction *objects.Interaction

		expectsErr string
		expects    *objects.InteractionResponse
	}{
		{
			name: "unset route",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID: "/a",
				}),
			},
			expectsErr: "modal path not found",
		},
		{
			name: "success",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationModalInteractionData{
					CustomID: "/a",
					Components: []*objects.InteractionResponseComponent{
						{
							Type:     objects.ComponentTypeInputText,
							CustomID: "abc",
							Value:    "def",
						},
						{
							Type:     objects.ComponentTypeInputText,
							CustomID: "abclol",
							Value:    "def",
						},
					},
				}),
			},
			init: func(t *testing.T, r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("this should not be called")
					},
					Function: func(ctx *ModalRouterCtx) error {
						assert.Equal(t, map[string]string{
							"abc": "def", "abclol": "def",
						}, ctx.ModalItems)
						ctx.SetContent("hello world!").Ephemeral()
						return nil
					},
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "hello world!",
					Flags:   64,
				},
			},
		},
		{
			name: "modal custom param",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationModalInteractionData{
					CustomID: "/a/hello",
					Components: []*objects.InteractionResponseComponent{
						{
							Type:     objects.ComponentTypeInputText,
							CustomID: "abc",
							Value:    "def",
						},
						{
							Type:     objects.ComponentTypeInputText,
							CustomID: "abclol",
							Value:    "def",
						},
					},
				}),
			},
			init: func(t *testing.T, r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a/:part",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("this should not be called")
					},
					Function: func(ctx *ModalRouterCtx) error {
						assert.Equal(t, map[string]string{
							"abc": "def", "abclol": "def",
						}, ctx.ModalItems)
						assert.Equal(t, map[string]string{"part": "hello"}, ctx.Params)
						ctx.SetContent("hello world!").Ephemeral()
						return nil
					},
				})
			},
			expects: &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "hello world!",
					Flags:   64,
				},
			},
		},
		{
			name: "modal error",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID: "/a",
				}),
			},
			init: func(_ *testing.T, r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("this should not be called")
					},
					Function: func(ctx *ModalRouterCtx) error {
						return errors.New("wumpus fled the scene")
					},
				})
			},
			expectsErr: "wumpus fled the scene",
		},
		{
			name: "modal panic",
			interaction: &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID: "/a",
				}),
			},
			init: func(_ *testing.T, r *ModalRouter) {
				r.AddModal(&ModalContent{
					Path: "/a",
					Contents: func(_ *ModalGenerationCtx) (string, []ModalContentItem) {
						panic("this should not be called")
					},
					Function: func(ctx *ModalRouterCtx) error {
						panic("wumpus fled the scene")
					},
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
			r := &ModalRouter{}
			if tt.init != nil {
				tt.init(t, r)
			}
			builtFunc := r.build(loaderPassthrough{
				rest:                  dummyRestClient,
				errHandler:            setError,
				modalRouter:           nil,
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
