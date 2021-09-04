package router

import (
	"reflect"
	"testing"

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

		respType objects.ResponseType
		data *objects.InteractionApplicationCommandCallbackData
		component bool
		globalAllowedMentions *objects.AllowedMentions

		wantErr   string
		wantResponse *objects.InteractionResponse
	}{
		{
			name:                  "no command response",
			component:             false,
			wantErr:               "expected data for command response",
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
			data:				   &objects.InteractionApplicationCommandCallbackData{},
			globalAllowedMentions: &objects.AllowedMentions{},
			wantResponse:          &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					AllowedMentions: &objects.AllowedMentions{},
				},
			},
		},
		{
			name:                  "default component data",
			component:             true,
			data:				   &objects.InteractionApplicationCommandCallbackData{},
			globalAllowedMentions: &objects.AllowedMentions{},
			wantResponse:          &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					AllowedMentions: &objects.AllowedMentions{},
				},
			},
		},
		{
			name:                  "override command type guess",
			component:             false,
			respType:			   69,
			data:				   &objects.InteractionApplicationCommandCallbackData{},
			wantResponse:          &objects.InteractionResponse{
				Type: 69,
				Data: &objects.InteractionApplicationCommandCallbackData{},
			},
		},
		{
			name:                  "override component type guess",
			component:             true,
			respType:			   69,
			data:				   &objects.InteractionApplicationCommandCallbackData{},
			wantResponse:          &objects.InteractionResponse{
				Type: 69,
				Data: &objects.InteractionApplicationCommandCallbackData{},
			},
		},
		{
			name: "full command",
			component:             false,
			data:				   &objects.InteractionApplicationCommandCallbackData{
				TTS:     true,
				Content: "testing testing 123",
				Embeds:  []*objects.Embed{
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
			wantResponse:          &objects.InteractionResponse{
				Type: objects.ResponseChannelMessageWithSource,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     true,
					Content: "testing testing 123",
					Embeds:  []*objects.Embed{
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
			name: "full component",
			component:             true,
			data:				   &objects.InteractionApplicationCommandCallbackData{
				TTS:     true,
				Content: "testing testing 123",
				Embeds:  []*objects.Embed{
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
			wantResponse:          &objects.InteractionResponse{
				Type: objects.ResponseUpdateMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     true,
					Content: "testing testing 123",
					Embeds:  []*objects.Embed{
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

func callBuilderFunction(t *testing.T, builder interface{}, funcName string, args ...interface{}) {
	t.Helper()
	r := reflect.ValueOf(builder).MethodByName(funcName)
	if r.IsZero() {
		t.Fatal("function does not exist")
	}
	reflectArgs := make([]reflect.Value, len(args))
	for i, v := range args {
		reflectArgs[i] = reflect.ValueOf(v)
	}
	if r.Kind() != reflect.Func {
		t.Fatal("not a function")
	}
	res := r.Call(reflectArgs)
	if len(res) != 1 {
		t.Fatal("arg count not correct for builder:", res)
	}
	if res[0].Interface() != builder {
		t.Fatal("the argument returned was not the builder")
	}
}

func TestComponentRouterCtx_SetEmbed(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "SetEmbed", &objects.Embed{Title: "a"})
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestCommandRouterCtx_SetEmbed(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "SetEmbed", &objects.Embed{Title: "a"})
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestComponentRouterCtx_AddEmbed(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "AddEmbed", &objects.Embed{Title: "a"})
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

func TestCommandRouterCtx_AddEmbed(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "AddEmbed", &objects.Embed{Title: "a"})
	assert.Equal(t, x.responseBuilder.ResponseData().Embeds, []*objects.Embed{{Title: "a"}})
}

var multipleComponentRowsRaw = []*objects.Component{
	{Type: objects.ComponentTypeActionRow, Components: []*objects.Component{{Label: "a"}}},
	{Type: objects.ComponentTypeActionRow, Components: []*objects.Component{{Label: "b"}}},
}

func TestComponentRouterCtx_AddComponentRow(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}})
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "b"}})
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_AddComponentRow(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}})
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "b"}})
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

var multipleComponentRows = [][]*objects.Component{
	{{Label: "a"}},
	{{Label: "b"}},
}

func TestComponentRouterCtx_SetComponentRows(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "c"}})
	callBuilderFunction(t, x, "SetComponentRows", multipleComponentRows)
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_SetComponentRows(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "c"}})
	callBuilderFunction(t, x, "SetComponentRows", multipleComponentRows)
	assert.Equal(t, x.responseBuilder.ResponseData().Components, multipleComponentRowsRaw)
}

func TestCommandRouterCtx_ClearComponents(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}})
	callBuilderFunction(t, x, "ClearComponents")
	assert.Equal(t, x.responseBuilder.ResponseData().Components, []*objects.Component{})
}

func TestComponentRouterCtx_ClearComponents(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "AddComponentRow", []*objects.Component{{Label: "a"}})
	callBuilderFunction(t, x, "ClearComponents")
	assert.Equal(t, x.responseBuilder.ResponseData().Components, []*objects.Component{})
}

func TestCommandRouterCtx_SetContent(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "SetContent", "a")
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "a")
}

func TestComponentRouterCtx_SetContent(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "SetContent", "a")
	assert.Equal(t, x.responseBuilder.ResponseData().Content, "a")
}

func TestCommandRouterCtx_SetAllowedMentions(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "SetAllowedMentions", &objects.AllowedMentions{Parse: []string{"a"}})
	assert.Equal(t, x.responseBuilder.ResponseData().AllowedMentions, &objects.AllowedMentions{Parse: []string{"a"}})
}

func TestComponentRouterCtx_SetAllowedMentions(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "SetAllowedMentions", &objects.AllowedMentions{Parse: []string{"a"}})
	assert.Equal(t, x.responseBuilder.ResponseData().AllowedMentions, &objects.AllowedMentions{Parse: []string{"a"}})
}

func TestCommandRouterCtx_SetTTS(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "SetTTS", true)
	assert.Equal(t, x.responseBuilder.ResponseData().TTS, true)
}

func TestComponentRouterCtx_SetTTS(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "SetTTS", true)
	assert.Equal(t, x.responseBuilder.ResponseData().TTS, true)
}

func TestCommandRouterCtx_Ephemeral(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "Ephemeral")
	assert.Equal(t, x.responseBuilder.ResponseData().Flags, 64)
}

func TestComponentRouterCtx_Ephemeral(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "Ephemeral")
	assert.Equal(t, x.responseBuilder.ResponseData().Flags, 64)
}

func TestCommandRouterCtx_ChannelMessageWithSource(t *testing.T) {
	x := &CommandRouterCtx{}
	callBuilderFunction(t, x, "ChannelMessageWithSource")
	assert.Equal(t, x.respType, objects.ResponseChannelMessageWithSource)
}

func TestComponentRouterCtx_ChannelMessageWithSource(t *testing.T) {
	x := &ComponentRouterCtx{}
	callBuilderFunction(t, x, "ChannelMessageWithSource")
	assert.Equal(t, x.respType, objects.ResponseChannelMessageWithSource)
}
