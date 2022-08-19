package router

import (
	"container/list"
	"context"
	"errors"
	"testing"

	"wumpgo.dev/wumpgo/objects"
	"github.com/stretchr/testify/assert"
)

func Test_findOption(t *testing.T) {
	tests := []struct {
		name string

		optName string
		options []*objects.ApplicationCommandOption

		wants *objects.ApplicationCommandOption
	}{
		{
			name: "nil slice",
		},
		{
			name:    "empty slice",
			options: []*objects.ApplicationCommandOption{},
		},
		{
			name:    "one option",
			optName: "opt1",
			options: []*objects.ApplicationCommandOption{
				{
					Name: "opt1",
				},
			},
			wants: &objects.ApplicationCommandOption{
				Name: "opt1",
			},
		},
		{
			name:    "two options",
			optName: "opt1",
			options: []*objects.ApplicationCommandOption{
				{
					Name: "opt1",
				},
				{
					Name: "opt2",
				},
			},
			wants: &objects.ApplicationCommandOption{
				Name: "opt1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wants := findOption(tt.optName, tt.options)
			assert.Equal(t, tt.wants, wants)
		})
	}
}

func TestCommand_mapOptions(t *testing.T) {
	tests := []struct {
		name string

		autocomplete bool
		data         *objects.ApplicationCommandInteractionData
		cmdOptions   []*objects.ApplicationCommandOption
		retOptions   []*objects.ApplicationCommandInteractionDataOption

		expectsErr string
		expects    map[string]any
	}{
		{
			name: "option not in command",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeInteger,
					Name:       "opt1",
				},
				{
					OptionType: objects.TypeString,
					Name:       "opt2",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeString,
					Name:  "opt3",
					Value: "123",
				},
			},
			expectsErr: "interaction option doesn't exist on command",
		},
		{
			name:         "autocomplete type mismatch",
			autocomplete: true,
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeInteger,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeString,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": "123",
			},
		},
		{
			name: "non-autocomplete type mismatch",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeInteger,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeString,
					Name:  "opt1",
					Value: "123",
				},
			},
			expectsErr: "mismatched interaction option",
		},
		{
			name: "channel option",
			data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeChannel,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeChannel,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": (ResolvableChannel)(resolvable[objects.Channel]{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
				}),
			},
		},
		{
			name: "role option",
			data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeRole,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeRole,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": (ResolvableRole)(resolvable[objects.Role]{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
				}),
			},
		},
		{
			name: "user option",
			data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeUser,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeUser,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": (ResolvableUser)(resolvableUser{resolvable[objects.User]{
					id:   "123",
					data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
				}}),
			},
		},
		{
			name: "string option",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeString,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeString,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": "123",
			},
		},
		{
			name: "int option",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeInteger,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeInteger,
					Name:  "opt1",
					Value: (float64)(69),
				},
			},
			expects: map[string]any{
				"opt1": 69,
			},
		},
		{
			name: "boolean option",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeBoolean,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeBoolean,
					Name:  "opt1",
					Value: true,
				},
			},
			expects: map[string]any{
				"opt1": true,
			},
		},
		{
			name: "mentionable option",
			data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeMentionable,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeMentionable,
					Name:  "opt1",
					Value: "123",
				},
			},
			expects: map[string]any{
				"opt1": (ResolvableMentionable)(resolvableMentionable{
					resolvable: resolvable[any]{
						id:   "123",
						data: &objects.ApplicationCommandInteractionData{DiscordBaseObject: objects.DiscordBaseObject{ID: 6921}},
					},
				}),
			},
		},
		{
			name: "double option",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeNumber,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeNumber,
					Name:  "opt1",
					Value: (float64)(69),
				},
			},
			expects: map[string]any{
				"opt1": (float64)(69),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the error handler to figure out if we got an error.
			var errResult error
			errHandler := func(err error) *objects.InteractionResponse {
				errResult = err
				return &objects.InteractionResponse{Type: 69}
			}

			// Call the function.
			c := Command{Options: tt.cmdOptions}
			resp, items := c.mapOptions(tt.autocomplete, tt.data, tt.retOptions, errHandler)
			if tt.expectsErr == "" {
				assert.NoError(t, errResult)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expects, items)
			} else {
				assert.EqualError(t, errResult, tt.expectsErr)
				assert.Nil(t, items)
				assert.Equal(t, &objects.InteractionResponse{Type: 69}, resp)
			}
		})
	}
}

var messageTargetData = &objects.ApplicationCommandInteractionData{
	TargetID: 1,
	Resolved: objects.ApplicationCommandInteractionDataResolved{
		Users: map[objects.Snowflake]objects.User{
			1: {
				DiscordBaseObject: objects.DiscordBaseObject{ID: 6921},
			},
		},
		Messages: map[objects.Snowflake]objects.Message{
			1: {
				DiscordBaseObject: objects.DiscordBaseObject{ID: 6921},
			},
		},
	},
}

var userTargetData = &objects.ApplicationCommandInteractionData{
	TargetID: 2,
	Resolved: objects.ApplicationCommandInteractionDataResolved{
		Users: map[objects.Snowflake]objects.User{
			2: {
				DiscordBaseObject: objects.DiscordBaseObject{ID: 6921},
			},
		},
		Messages: map[objects.Snowflake]objects.Message{
			1: {
				DiscordBaseObject: objects.DiscordBaseObject{ID: 6921},
			},
		},
	},
}

func TestCommand_execute(t *testing.T) {
	tests := []struct {
		name string

		paramsCheck  func(*testing.T, map[string]any)
		data         *objects.ApplicationCommandInteractionData
		retOptions   []*objects.ApplicationCommandInteractionDataOption
		cmdOptions   []*objects.ApplicationCommandOption
		cmdAllowed   *objects.AllowedMentions
		noMiddleware bool
		noFunc       bool

		throwsErr  string
		panic      bool
		expectsErr string
	}{
		{
			name: "message target",
			data: messageTargetData,
			paramsCheck: func(t *testing.T, m map[string]any) {
				t.Helper()
				assert.Equal(t, resolvable[objects.Message]{
					id:   "1",
					data: messageTargetData,
				}, m["/target"])
			},
		},
		{
			name: "user target",
			data: userTargetData,
			paramsCheck: func(t *testing.T, m map[string]any) {
				t.Helper()
				assert.Equal(t, (ResolvableUser)(resolvableUser{resolvable[objects.User]{
					id:   "2",
					data: userTargetData,
				}}), m["/target"])
			},
		},
		{
			name: "failed to map options",
			cmdOptions: []*objects.ApplicationCommandOption{
				{
					OptionType: objects.TypeString,
					Name:       "opt1",
				},
			},
			retOptions: []*objects.ApplicationCommandInteractionDataOption{
				{
					Type:  objects.TypeInteger,
					Name:  "opt1",
					Value: 123,
				},
			},
			expectsErr: "mismatched interaction option",
		},
		{
			name:       "returned error",
			throwsErr:  "cat tripped on wire",
			expectsErr: "cat tripped on wire",
		},
		{
			name:       "panic",
			throwsErr:  "cat tripped on wire",
			panic:      true,
			expectsErr: "cat tripped on wire",
		},
		{
			name:       "command allowed mentions",
			cmdAllowed: &objects.AllowedMentions{Parse: []string{"a", "b", "c"}},
		},
		{
			name: "success",
		},
		{
			name:         "no middleware success",
			noMiddleware: true,
		},
		{
			name:         "no middleware error",
			noMiddleware: true,
			throwsErr:    "cat tripped on wire",
			expectsErr:   "cat tripped on wire",
		},
		{
			name:       "no handler",
			noFunc:     true,
			expectsErr: "expected data for command response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure data isn't nil.
			if tt.data == nil {
				tt.data = &objects.ApplicationCommandInteractionData{Options: tt.retOptions}
			}

			// Defines a list of middleware.
			l := list.New()
			if !tt.noMiddleware {
				l.PushBack((MiddlewareFunc)(middleware1))
				l.PushBack((MiddlewareFunc)(middleware2))
			}

			// Defines the command.
			c := &Command{Options: tt.cmdOptions, AllowedMentions: tt.cmdAllowed}

			// Defines the dummy interaction.
			dummyInteraction := &objects.Interaction{}

			// Defines the command function.
			f := func(ctx *CommandRouterCtx) error {
				// Defines this as a helper function.
				t.Helper()

				// Check the passed through items are what was expected.
				assert.Equal(t, c, ctx.Command)
				assert.Same(t, dummyInteraction, ctx.Interaction)
				assert.Same(t, dummyRestClient, ctx.RESTClient)
				assert.Same(t, tt.cmdAllowed, ctx.globalAllowedMentions)

				// Check the 2 expected middlewares were ran.
				if !tt.noMiddleware {
					assert.Equal(t, "middleware1", ctx.Options["middleware1"])
					assert.Equal(t, 2, ctx.Options["middleware2"])
				}

				// Check the params.
				if tt.paramsCheck != nil {
					tt.paramsCheck(t, ctx.Options)
				}

				// Check if we are supposed to return an error.
				if tt.throwsErr != "" {
					if tt.panic {
						panic(tt.throwsErr)
					}
					return errors.New(tt.throwsErr)
				}

				// Return no errors.
				ctx.SetContent("hello world")
				return nil
			}
			if !tt.noFunc {
				c.Function = f
			}

			// Set the error handler to figure out if we got an error.
			var errResult error
			errHandler := func(err error) *objects.InteractionResponse {
				errResult = err
				return &objects.InteractionResponse{Type: 69}
			}

			// Execute the command.
			resp := c.execute(context.Background(), commandExecutionOptions{
				restClient:       dummyRestClient,
				exceptionHandler: errHandler,
				interaction:      dummyInteraction,
				data:             tt.data,
				options:          tt.retOptions,
			}, l)
			if tt.expectsErr == "" {
				assert.NoError(t, errResult)
				if tt.noFunc {
					assert.Nil(t, resp)
				} else {
					assert.Equal(t, &objects.InteractionResponse{
						Type: 4,
						Data: &objects.InteractionApplicationCommandCallbackData{
							TTS:             false,
							Content:         "hello world",
							AllowedMentions: tt.cmdAllowed,
						},
					}, resp)
				}
			} else {
				assert.EqualError(t, errResult, tt.expectsErr)
				assert.Equal(t, &objects.InteractionResponse{Type: 69}, resp)
			}
		})
	}
}

func TestCommand_Groups(t *testing.T) {
	g1 := &CommandGroup{}
	g2 := &CommandGroup{parent: g1}
	g3 := &CommandGroup{parent: g2}
	t.Run("groups", func(t *testing.T) {
		cmd := &Command{parent: g3}
		assert.Equal(t, []*CommandGroup{g1, g2, g3}, cmd.Groups())
	})
	t.Run("no groups", func(t *testing.T) {
		cmd := &Command{}
		assert.Equal(t, []*CommandGroup{}, cmd.Groups())
	})
}
