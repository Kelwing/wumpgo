package router

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func Test_responseBuilder_ResponseData(t *testing.T) {
	tests := []struct {
		name string

		data *objects.InteractionApplicationCommandCallbackData
	}{
		{
			name: "nil data",
		},
		{
			name: "not nil data",
			data: &objects.InteractionApplicationCommandCallbackData{TTS: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{dataPtr: tt.data}
			respData := b.ResponseData()
			if tt.data == nil {
				assert.Equal(t, &objects.InteractionApplicationCommandCallbackData{}, respData)
			} else {
				assert.Equal(t, tt.data, respData)
			}
		})
	}
}

func Test_responseBuilder_buildResponse(t *testing.T) {
	tests := []struct {
		name string

		respType              objects.ResponseType
		data                  *objects.InteractionApplicationCommandCallbackData
		component             bool
		globalAllowedMentions *objects.AllowedMentions

		wantErr      string
		wantResponse *objects.InteractionResponse
	}{
		{
			name:      "no command response",
			component: false,
			wantErr:   "expected data for command response",
		},
		{
			name:                  "no component response",
			component:             true,
			globalAllowedMentions: &objects.AllowedMentions{},
			wantResponse:          &objects.InteractionResponse{Type: objects.ResponseDeferredMessageUpdate},
		},
		{
			name:                  "default command data",
			component:             false,
			data:                  &objects.InteractionApplicationCommandCallbackData{},
			globalAllowedMentions: &objects.AllowedMentions{},
			wantResponse: &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					AllowedMentions: &objects.AllowedMentions{},
				},
			},
		},
		{
			name:                  "default component data",
			component:             true,
			data:                  &objects.InteractionApplicationCommandCallbackData{},
			globalAllowedMentions: &objects.AllowedMentions{},
			wantResponse: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					AllowedMentions: &objects.AllowedMentions{},
				},
			},
		},
		{
			name:      "override command type guess",
			component: false,
			respType:  69,
			data:      &objects.InteractionApplicationCommandCallbackData{},
			wantResponse: &objects.InteractionResponse{
				Type: 69,
				Data: &objects.InteractionApplicationCommandCallbackData{},
			},
		},
		{
			name:      "override component type guess",
			component: true,
			respType:  69,
			data:      &objects.InteractionApplicationCommandCallbackData{},
			wantResponse: &objects.InteractionResponse{
				Type: 69,
				Data: &objects.InteractionApplicationCommandCallbackData{},
			},
		},
		{
			name:      "full command",
			component: false,
			data: &objects.InteractionApplicationCommandCallbackData{
				TTS:     true,
				Content: "testing testing 123",
				Embeds: []*objects.Embed{
					{
						Title: "hello world",
					},
				},
				AllowedMentions: &objects.AllowedMentions{
					Parse:       []string{"1"},
					Roles:       []objects.Snowflake{2},
					Users:       []objects.Snowflake{3},
					RepliedUser: true,
				},
				Components: []*objects.Component{
					{
						Label: "test",
					},
				},
			},
			wantResponse: &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     true,
					Content: "testing testing 123",
					Embeds: []*objects.Embed{
						{
							Title: "hello world",
						},
					},
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"1"},
						Roles:       []objects.Snowflake{2},
						Users:       []objects.Snowflake{3},
						RepliedUser: true,
					},
					Components: []*objects.Component{
						{
							Label: "test",
						},
					},
				},
			},
		},
		{
			name:      "full component",
			component: true,
			data: &objects.InteractionApplicationCommandCallbackData{
				TTS:     true,
				Content: "testing testing 123",
				Embeds: []*objects.Embed{
					{
						Title: "hello world",
					},
				},
				AllowedMentions: &objects.AllowedMentions{
					Parse:       []string{"1"},
					Roles:       []objects.Snowflake{2},
					Users:       []objects.Snowflake{3},
					RepliedUser: true,
				},
				Components: []*objects.Component{
					{
						Label: "test",
					},
				},
			},
			wantResponse: &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     true,
					Content: "testing testing 123",
					Embeds: []*objects.Embed{
						{
							Title: "hello world",
						},
					},
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"1"},
						Roles:       []objects.Snowflake{2},
						Users:       []objects.Snowflake{3},
						RepliedUser: true,
					},
					Components: []*objects.Component{
						{
							Label: "test",
						},
					},
				},
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

			// Create the response builder and call the builder.
			b := responseBuilder{respType: tt.respType, dataPtr: tt.data}
			response := b.buildResponse(tt.component, setError, tt.globalAllowedMentions)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantResponse, response)
		})
	}
}

func Test_responseBuilder_editEmbed(t *testing.T) {
	tests := []struct {
		name string

		initEmbeds []*objects.Embed
		embed      *objects.Embed
		append     bool
		expected   []*objects.Embed
	}{
		{
			name:     "all nil",
			expected: ([]*objects.Embed)(nil),
		},
		{
			name:       "non append",
			initEmbeds: []*objects.Embed{{}},
			embed:      &objects.Embed{Description: "a"},
			append:     false,
			expected:   []*objects.Embed{{Description: "a"}},
		},
		{
			name:       "append",
			initEmbeds: []*objects.Embed{{}},
			embed:      &objects.Embed{Description: "a"},
			append:     true,
			expected:   []*objects.Embed{{}, {Description: "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{}
			if tt.initEmbeds != nil {
				b.ResponseData().Embeds = tt.initEmbeds
			}
			b.editEmbed(tt.embed, tt.append)
			assert.Equal(t, tt.expected, b.ResponseData().Embeds)
		})
	}
}

func Test_responseBuilder_editComponent(t *testing.T) {
	tests := []struct {
		name string

		initComponents []*objects.Component
		component      *objects.Component
		append         bool
		expected       []*objects.Component
	}{
		{
			name:     "all nil",
			expected: ([]*objects.Component)(nil),
		},
		{
			name:           "non append",
			initComponents: []*objects.Component{{}},
			component:      &objects.Component{Placeholder: "test"},
			append:         false,
			expected:       []*objects.Component{{Placeholder: "test"}},
		},
		{
			name:           "append",
			initComponents: []*objects.Component{{}},
			component:      &objects.Component{Placeholder: "a"},
			append:         true,
			expected:       []*objects.Component{{}, {Placeholder: "a"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := responseBuilder{}
			if tt.initComponents != nil {
				b.ResponseData().Components = tt.initComponents
			}
			b.editComponent(tt.component, tt.append)
			assert.Equal(t, tt.expected, b.ResponseData().Components)
		})
	}
}

func TestComponentRouterCtx_SetEmbed(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestCommandRouterCtx_SetEmbed(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestModalRouterCtx_SetEmbed(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestComponentRouterCtx_AddEmbed(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestCommandRouterCtx_AddEmbed(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestModalRouterCtx_AddEmbed(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddEmbed", &objects.Embed{Title: "a"}))
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

var multipleComponentRowsRaw = []*objects.Component{
	{Type: objects.ComponentTypeActionRow, Components: []*objects.Component{{Label: "a"}}},
	{Type: objects.ComponentTypeActionRow, Components: []*objects.Component{{Label: "b"}}},
}

func TestComponentRouterCtx_AddComponentRow(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "b"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_AddComponentRow(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "b"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestModalRouterCtx_AddComponentRow(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "b"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

var multipleComponentRows = [][]*objects.Component{
	{{Label: "a"}},
	{{Label: "b"}},
}

func TestComponentRouterCtx_SetComponentRows(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "c"}}))
	assert.NoError(t, callBuilderFunction(t, x, "SetComponentRows", multipleComponentRows))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_SetComponentRows(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "c"}}))
	assert.NoError(t, callBuilderFunction(t, x, "SetComponentRows", multipleComponentRows))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestModalRouterCtx_SetComponentRows(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "c"}}))
	assert.NoError(t, callBuilderFunction(t, x, "SetComponentRows", multipleComponentRows))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_ClearComponents(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "ClearComponents"))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, []*objects.Component{})
}

func TestComponentRouterCtx_ClearComponents(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "ClearComponents"))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, []*objects.Component{})
}

func TestModalRouterCtx_ClearComponents(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}}))
	assert.NoError(t, callBuilderFunction(t, x, "ClearComponents"))
	assert.Equal(t, x.responseBuilder.ResponseData().Components, []*objects.Component{})
}

func TestCommandRouterCtx_SetContent(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContent", "a"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "a")
}

func TestComponentRouterCtx_SetContent(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContent", "a"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "a")
}

func TestModalRouterCtx_SetContent(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContent", "a"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "a")
}

func TestCommandRouterCtx_SetContentf(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContentf", "hello %s", "world"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "hello world")
}

func TestComponentRouterCtx_SetContentf(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContentf", "hello %s", "world"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "hello world")
}

func TestModalRouterCtx_SetContentf(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetContentf", "hello %s", "world"))
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "hello world")
}

func TestCommandRouterCtx_SetAllowedMentions(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetAllowedMentions", &objects.AllowedMentions{Parse: []string{"a"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().AllowedMentions, &objects.AllowedMentions{Parse: []string{"a"}})
}

func TestComponentRouterCtx_SetAllowedMentions(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetAllowedMentions", &objects.AllowedMentions{Parse: []string{"a"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().AllowedMentions, &objects.AllowedMentions{Parse: []string{"a"}})
}

func TestModalRouterCtx_SetAllowedMentions(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetAllowedMentions", &objects.AllowedMentions{Parse: []string{"a"}}))
	assert.Equal(t, x.responseBuilder.ResponseData().AllowedMentions, &objects.AllowedMentions{Parse: []string{"a"}})
}

func TestCommandRouterCtx_SetTTS(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetTTS", true))
	assert.Equal(t, x.responseBuilder.ResponseData().TTS, true)
}

func TestComponentRouterCtx_SetTTS(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetTTS", true))
	assert.Equal(t, x.responseBuilder.ResponseData().TTS, true)
}

func TestModalRouterCtx_SetTTS(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "SetTTS", true))
	assert.Equal(t, x.responseBuilder.ResponseData().TTS, true)
}

func TestCommandRouterCtx_Ephemeral(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "Ephemeral"))
	assert.Equal(t, (objects.MessageFlag)(64), x.responseBuilder.ResponseData().Flags)
}

func TestComponentRouterCtx_Ephemeral(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "Ephemeral"))
	assert.Equal(t, (objects.MessageFlag)(64), x.responseBuilder.ResponseData().Flags)
}

func TestModalRouterCtx_Ephemeral(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "Ephemeral"))
	assert.Equal(t, (objects.MessageFlag)(64), x.responseBuilder.ResponseData().Flags)
}

func TestCommandRouterCtx_AttachBytes(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AttachBytes", []byte("a"), "file.txt", ""))
	assert.Equal(t, "file.txt", x.responseBuilder.ResponseData().Files[0].Filename)
}

func TestCommandRouterCtx_AttachFile(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AttachFile", &objects.DiscordFile{
		Buffer:   bytes.NewBuffer([]byte("a")),
		Filename: "file.txt",
	}))
	assert.Equal(t, "file.txt", x.responseBuilder.ResponseData().Files[0].Filename)
}

func TestComponentRouterCtx_AttachBytes(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AttachBytes", []byte("a"), "file.txt", ""))
	assert.Equal(t, "file.txt", x.responseBuilder.ResponseData().Files[0].Filename)
}

func TestComponentRouterCtx_AttachFile(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "AttachFile", &objects.DiscordFile{
		Buffer:   bytes.NewBuffer([]byte("a")),
		Filename: "file.txt",
	}))
	assert.Equal(t, "file.txt", x.responseBuilder.ResponseData().Files[0].Filename)
}

func TestCommandRouterCtx_ChannelMessageWithSource(t *testing.T) {
	x := &CommandRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "ChannelMessageWithSource"))
	assert.Equal(t, x.respType, objects.ResponseChannelMessageWithSource)
}

func TestComponentRouterCtx_ChannelMessageWithSource(t *testing.T) {
	x := &ComponentRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "ChannelMessageWithSource"))
	assert.Equal(t, x.respType, objects.ResponseChannelMessageWithSource)
}

func TestModalRouterCtx_ChannelMessageWithSource(t *testing.T) {
	x := &ModalRouterCtx{}
	assert.NoError(t, callBuilderFunction(t, x, "ChannelMessageWithSource"))
	assert.Equal(t, x.respType, objects.ResponseChannelMessageWithSource)
}

func buildWithModalRouter(r any, mount bool, setError func(e error) *objects.InteractionResponse, app HandlerAccepter) {
	modalRouter := &ModalRouter{}
	modalRouter.AddModal(&ModalContent{
		Path: "/test/:content",
		Contents: func(ctx *ModalGenerationCtx) (name string, contents []ModalContentItem) {
			return ctx.Path, []ModalContentItem{
				{
					Label: "test",
				},
			}

		},
		Function: func(ctx *ModalRouterCtx) error {
			panic("this should not be called!")
		},
	})
	rb := RouterLoader().ErrorHandler(setError)
	if mount {
		rb.ModalRouter(modalRouter)
	}
	switch x := r.(type) {
	case *CommandRouter:
		rb.CommandRouter(x)
	case *ComponentRouter:
		rb.ComponentRouter(x)
	}
	rb.Build(app)
}

func TestComponentRouterCtx_WithModalPath(t *testing.T) {
	tests := []struct {
		name string

		path              string
		mountModalHandler bool
		selectMenu        bool

		expectsErr string
		expects    *objects.InteractionResponse
	}{
		// Button

		{
			name:       "button no modal router",
			expectsErr: "modal router is unset",
		},
		{
			name:              "button unknown path",
			path:              "/unknown",
			expectsErr:        "modal path not found",
			mountModalHandler: true,
		},
		{
			name:              "button valid path",
			path:              "/test/test",
			mountModalHandler: true,
			expects: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:  objects.ComponentTypeInputText,
									Label: "test",
									Style: objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/test/test",
					Title:    "/test/test",
				},
			},
		},

		// Select menu

		{
			name:       "select menu no modal router",
			expectsErr: "modal router is unset",
			selectMenu: true,
		},
		{
			name:              "select menu unknown path",
			path:              "/unknown",
			expectsErr:        "modal path not found",
			mountModalHandler: true,
			selectMenu:        true,
		},
		{
			name:              "select menu valid path",
			path:              "/test/test",
			mountModalHandler: true,
			selectMenu:        true,
			expects: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:  objects.ComponentTypeInputText,
									Label: "test",
									Style: objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/test/test",
					Title:    "/test/test",
				},
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

			// Create the router and mock route.
			r := &ComponentRouter{}
			if tt.selectMenu {
				r.RegisterSelectMenu("/test", func(ctx *ComponentRouterCtx, _ []string) error {
					t.Helper()
					return ctx.WithModalPath(tt.path)
				})
			} else {
				r.RegisterButton("/test", func(ctx *ComponentRouterCtx) error {
					t.Helper()
					return ctx.WithModalPath(tt.path)
				})
			}

			// Build the router.
			f := &fakeBuildHandlerAccepter{}
			buildWithModalRouter(r, tt.mountModalHandler, setError, f)
			require.NotNil(t, f.componentHandler)

			// Call the component handler.
			x := objects.ComponentTypeButton
			if tt.selectMenu {
				x = objects.ComponentTypeSelectMenu
			}
			resp := f.componentHandler(context.Background(), &objects.Interaction{
				Data: jsonify(t, objects.ApplicationComponentInteractionData{
					CustomID:      "/test",
					ComponentType: x,
				}),
			})

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

func TestCommandRouterCtx_WithModalPath(t *testing.T) {
	tests := []struct {
		name string

		path              string
		mountModalHandler bool

		expectsErr string
		expects    *objects.InteractionResponse
	}{
		{
			name:       "no modal router",
			expectsErr: "modal router is unset",
		},
		{
			name:              "unknown path",
			path:              "/unknown",
			expectsErr:        "modal path not found",
			mountModalHandler: true,
		},
		{
			name:              "valid path",
			path:              "/test/test",
			mountModalHandler: true,
			expects: &objects.InteractionResponse{
				Type: objects.ResponseModal,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Components: []*objects.Component{
						{
							Type: objects.ComponentTypeActionRow,
							Components: []*objects.Component{
								{
									Type:  objects.ComponentTypeInputText,
									Label: "test",
									Style: objects.ButtonStyle(objects.TextStyleParagraph),
								},
							},
						},
					},
					CustomID: "/test/test",
					Title:    "/test/test",
				},
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

			// Create the router and mock command.
			r := &CommandRouter{}
			_, err = r.NewCommandBuilder("test").Handler(func(ctx *CommandRouterCtx) error {
				t.Helper()
				return ctx.WithModalPath(tt.path)
			}).Build()
			require.NoError(t, err)

			// Build the router.
			f := &fakeBuildHandlerAccepter{}
			buildWithModalRouter(r, tt.mountModalHandler, setError, f)
			require.NotNil(t, f.commandHandler)

			// Call the command handler.
			resp := f.commandHandler(context.Background(), &objects.Interaction{
				Data: jsonify(t, objects.ApplicationCommandInteractionData{
					DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
					Name:              "test",
					Type:              objects.CommandTypeChatInput,
					Version:           1,
					Resolved:          objects.ApplicationCommandInteractionDataResolved{},
				}),
			})

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
