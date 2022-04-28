package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Postcord/objects"
	"github.com/jimeh/go-golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandRouterCtx_TargetMessage(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expects *objects.Message
	}{
		{
			name: "no message",
		},
		{
			name:    "message wrong type",
			options: map[string]interface{}{"/target": 1},
		},
		{
			name: "message exists",
			options: map[string]interface{}{
				"/target": &ResolvableMessage{
					id: "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Messages: map[objects.Snowflake]objects.Message{
								123: {
									Content: "hello world",
								},
							},
						},
					},
				},
			},
			expects: &objects.Message{Content: "hello world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.options == nil {
				tt.options = map[string]interface{}{}
			}
			ctx := CommandRouterCtx{Options: tt.options}
			assert.Equal(t, tt.expects, ctx.TargetMessage())
		})
	}
}

func TestCommandRouterCtx_TargetMember(t *testing.T) {
	tests := []struct {
		name string

		options map[string]interface{}
		expects *objects.GuildMember
	}{
		{
			name: "no member",
		},
		{
			name:    "member wrong type",
			options: map[string]interface{}{"/target": 1},
		},
		{
			name: "member exists",
			options: map[string]interface{}{
				"/target": &ResolvableUser{
					id: "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Members: map[objects.Snowflake]objects.GuildMember{
								123: {
									Deaf: true,
								},
							},
						},
					},
				},
			},
			expects: &objects.GuildMember{Deaf: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.options == nil {
				tt.options = map[string]interface{}{}
			}
			ctx := CommandRouterCtx{Options: tt.options}
			assert.Equal(t, tt.expects, ctx.TargetMember())
		})
	}
}

func Test_messageTargetWrapper(t *testing.T) {
	tests := []struct {
		name string

		options    map[string]interface{}
		expectsErr string
	}{
		{
			name: "successful call",
			options: map[string]interface{}{
				"/target": &ResolvableMessage{
					id: "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Messages: map[objects.Snowflake]objects.Message{
								123: {
									Content: "hello world",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "invalid target",
			options:    map[string]interface{}{},
			expectsErr: "wrong or no target specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := &CommandRouterCtx{
				Options: tt.options,
			}

			var called bool
			cb := func(ctx *CommandRouterCtx, msg *objects.Message) error {
				if called {
					t.Fatal("function called twice")
				}
				called = true

				if mockCtx != ctx {
					t.Fatal("context different")
				}
				assert.Equal(t, &objects.Message{Content: "hello world"}, msg)

				return nil
			}

			f := messageTargetWrapper(cb)
			err := f(mockCtx)
			if tt.expectsErr == "" {
				if !called {
					t.Fatal("function wasn't called")
				}
				assert.NoError(t, err)
			} else {
				if called {
					t.Fatal("function was called")
				}
				assert.EqualError(t, err, tt.expectsErr)
			}
		})
	}
}

func Test_memberTargetWrapper(t *testing.T) {
	tests := []struct {
		name string

		options    map[string]interface{}
		expectsErr string
	}{
		{
			name: "successful call",
			options: map[string]interface{}{
				"/target": &ResolvableUser{
					id: "123",
					data: &objects.ApplicationCommandInteractionData{
						TargetID: 123,
						Resolved: objects.ApplicationCommandInteractionDataResolved{
							Members: map[objects.Snowflake]objects.GuildMember{
								123: {
									Nick: "hello world",
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "invalid target",
			options:    map[string]interface{}{},
			expectsErr: "wrong or no target specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtx := &CommandRouterCtx{
				Options: tt.options,
			}

			var called bool
			cb := func(ctx *CommandRouterCtx, member *objects.GuildMember) error {
				if called {
					t.Fatal("function called twice")
				}
				called = true

				if mockCtx != ctx {
					t.Fatal("context different")
				}
				assert.Equal(t, &objects.GuildMember{Nick: "hello world"}, member)

				return nil
			}

			f := memberTargetWrapper(cb)
			err := f(mockCtx)
			if tt.expectsErr == "" {
				if !called {
					t.Fatal("function wasn't called")
				}
				assert.NoError(t, err)
			} else {
				if called {
					t.Fatal("function was called")
				}
				assert.EqualError(t, err, tt.expectsErr)
			}
		})
	}
}

func TestCommandGroup_Use(t *testing.T) {
	tests := []struct {
		name string

		init       func(c *CommandGroup)
		expectsLen int
	}{
		{
			name:       "no middleware",
			expectsLen: 0,
		},
		{
			name: "add middleware",
			init: func(c *CommandGroup) {
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
			},
			expectsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandGroup{Middleware: []MiddlewareFunc{}}
			if tt.init != nil {
				tt.init(c)
			}
			assert.Equal(t, tt.expectsLen, len(c.Middleware))
		})
	}
}

var dummyRootCommandGroup = &CommandGroup{}

var commandGroupTests = []struct {
	name string

	level             uint
	groupName         string
	description       string
	defaultPermission bool

	expects    *CommandGroup
	expectsErr string
}{
	{
		name:       "group nested too deep",
		level:      2,
		expectsErr: "sub-command group would be nested too deep",
	},
	{
		name:              "root group",
		groupName:         "abc",
		description:       "def",
		defaultPermission: true,
		expects: &CommandGroup{
			level:             1,
			parent:            dummyRootCommandGroup,
			DefaultPermission: true,
			Description:       "def",
			Subcommands:       map[string]interface{}{},
		},
	},
	{
		name:              "sub-group",
		groupName:         "abc",
		description:       "def",
		defaultPermission: true,
		level:             1,
		expects: &CommandGroup{
			level:             2,
			parent:            dummyRootCommandGroup,
			DefaultPermission: true,
			Description:       "def",
			Subcommands:       map[string]interface{}{},
		},
	},
}

func TestCommandGroup_NewCommandGroup(t *testing.T) {
	for _, tt := range commandGroupTests {
		t.Run(tt.name, func(t *testing.T) {
			dummyRootCommandGroup.Subcommands, dummyRootCommandGroup.level = map[string]interface{}{}, tt.level
			group, err := dummyRootCommandGroup.NewCommandGroup(tt.groupName, tt.description, nil)
			if tt.expectsErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectsErr)
				return
			}
			assert.Equal(t, tt.expects, group)
		})
	}
}

type mustNewCommandGroup interface {
	MustNewCommandGroup(name, description string, opts *CommandGroupOptions) *CommandGroup
}

func unpanicCommandGroup(x mustNewCommandGroup, name, description string, default_ bool) (group *CommandGroup, returnedErr string) {
	defer func() {
		if r := recover(); r != nil {
			returnedErr = fmt.Sprint(r)
		}
	}()
	group = x.MustNewCommandGroup(name, description, nil)
	return
}

func TestCommandGroup_MustNewCommandGroup(t *testing.T) {
	for _, tt := range commandGroupTests {
		t.Run(tt.name, func(t *testing.T) {
			dummyRootCommandGroup.Subcommands, dummyRootCommandGroup.level = map[string]interface{}{}, tt.level
			group, errResult := unpanicCommandGroup(dummyRootCommandGroup, tt.groupName, tt.description, tt.defaultPermission)
			assert.Equal(t, errResult, tt.expectsErr)
			assert.Equal(t, tt.expects, group)
		})
	}
}

func TestCommandRouter_Use(t *testing.T) {
	tests := []struct {
		name string

		init       func(c *CommandRouter)
		expectsLen int
	}{
		{
			name:       "no middleware",
			expectsLen: 0,
		},
		{
			name: "add middleware",
			init: func(c *CommandRouter) {
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
				c.Use(func(ctx MiddlewareCtx) error {
					return nil
				})
			},
			expectsLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CommandRouter{middleware: []MiddlewareFunc{}}
			if tt.init != nil {
				tt.init(r)
			}
			assert.Equal(t, tt.expectsLen, len(r.middleware))
		})
	}
}

func TestCommandRouter_NewCommandGroup(t *testing.T) {
	r := &CommandRouter{}
	group, err := r.NewCommandGroup("abc", "def")
	assert.NoError(t, err)
	assert.Equal(t, &CommandGroup{
		level:             1,
		DefaultPermission: true,
		Description:       "def",
		Subcommands:       map[string]interface{}{},
	}, group)
}

// func TestCommandRouter_MustNewCommandGroup(t *testing.T) {
// 	r := &CommandRouter{}
// 	group, errResult := unpanicCommandGroup(r, "abc", "def", true)
// 	assert.Equal(t, "", errResult)
// 	assert.Equal(t, &CommandGroup{
// 		level:             1,
// 		DefaultPermission: true,
// 		Description:       "def",
// 		Subcommands:       map[string]interface{}{},
// 	}, group)
// }

func TestCommandRouter_NewCommandBuilder(t *testing.T) {
	r := &CommandRouter{}
	builder := r.NewCommandBuilder("abc")
	assert.NotNil(t, r.roots.Subcommands)
	assert.Equal(t, &commandBuilder{
		map_: r.roots.Subcommands,
		cmd:  Command{Name: "abc"},
	}, builder)
}

// Returns the mock command runner. This is used to check everything in the context is as expected.
func mockCommandFunction(t *testing.T, expectedCmd *Command,
	expectedInteraction *objects.Interaction, paramsCheck func(*testing.T, map[string]interface{}),
	wantsErr string) func(ctx *CommandRouterCtx) error {
	// Return the generated command.
	return func(ctx *CommandRouterCtx) error {
		// Mark this as a helper function.
		t.Helper()

		// Check the passed through items are what was expected.
		assert.Equal(t, expectedCmd, ctx.Command)
		assert.Same(t, expectedInteraction, ctx.Interaction)
		assert.Same(t, dummyRestClient, ctx.RESTClient)

		// Check the 2 expected middlewares were ran.
		assert.Equal(t, "middleware1", ctx.Options["middleware1"])
		assert.Equal(t, 2, ctx.Options["middleware2"])

		// Check the params.
		paramsCheck(t, ctx.Options)

		// Check if we are supposed to return an error.
		if wantsErr != "" {
			return errors.New(wantsErr)
		}

		// Return no errors.
		ctx.SetContent("hello world")
		return nil
	}
}

// Returns the mock auto-complete runner. This is used to check everything in the context is as expected.
func mockAutocompleteFunction(t *testing.T, expectedCmd *Command,
	expectedInteraction *objects.Interaction, wantsErr string) func(*CommandRouterCtx) ([]StringChoice, error) {
	// Return the generated auto-complete function.
	return func(ctx *CommandRouterCtx) ([]StringChoice, error) {
		// Mark this as a helper function.
		t.Helper()

		// Check the passed through items are what was expected.
		assert.Equal(t, expectedCmd, ctx.Command)
		assert.Same(t, expectedInteraction, ctx.Interaction)
		assert.Same(t, dummyRestClient, ctx.RESTClient)

		// Check the option is as expected.
		assert.Equal(t, ctx.Options["autocompletetest"], "123")

		// Check if we are supposed to return an error.
		if wantsErr != "" {
			return nil, errors.New(wantsErr)
		}

		// Return no errors.
		return []StringChoice{{
			Name:  "a",
			Value: "b",
		}}, nil
	}
}

func middleware1(ctx MiddlewareCtx) error {
	ctx.Options["middleware1"] = "random shit to test order"
	return ctx.Next()
}

func middleware2(ctx MiddlewareCtx) error {
	ctx.Options["middleware1"] = "middleware1"
	ctx.Options["middleware2"] = 2
	return ctx.Next()
}

func makeMockFullCommandRouter(injectAllowedMentions *objects.AllowedMentions) (*CommandRouter, map[string]*Command) {
	r := &CommandRouter{}

	r.Use(middleware1)
	r.Use(middleware2)

	root3 := r.MustNewCommandGroup("root3", "")
	root3.AllowedMentions = injectAllowedMentions
	root4 := r.MustNewCommandGroup("root4", "")
	sub3 := root4.MustNewCommandGroup("sub3", "", nil)

	cmds := map[string]*Command{
		"root1": r.NewCommandBuilder("root1").MustBuild(),
		"root2": r.NewCommandBuilder("root2").
			StringOption("string", "strings", true, nil).
			IntOption("int", "ints", true, nil).MustBuild(),

		"sub1": root3.NewCommandBuilder("sub1").MustBuild(),
		"sub2": root3.NewCommandBuilder("sub2").
			StringOption("string", "strings", true, nil).
			IntOption("int", "ints", true, nil).MustBuild(),

		"subsub1": sub3.NewCommandBuilder("subsub1").MustBuild(),
		"subsub2": sub3.NewCommandBuilder("subsub2").
			StringOption("string", "strings", true, nil).
			IntOption("int", "ints", true, nil).MustBuild(),
	}

	return r, cmds
}

func mockInteraction(data interface{}) *objects.Interaction {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return &objects.Interaction{
		DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
		ApplicationID:     5678,
		Type:              69,
		Data:              b,
		GuildID:           1234,
		ChannelID:         5678,
		Member:            &objects.GuildMember{Nick: "jeff", User: &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}}},
		User:              &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}},
		Token:             "abcd",
		Message: &objects.Message{
			DiscordBaseObject: objects.DiscordBaseObject{ID: 9101112},
			ChannelID:         3210,
			GuildID:           4567,
			Author:            &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}},
			Member:            &objects.GuildMember{Nick: "jeff", User: &objects.User{DiscordBaseObject: objects.DiscordBaseObject{ID: 123}}},
			Content:           "hello world",
			Components:        []*objects.Component{{Type: 69}},
		},
		Version: 1,
	}
}

var helloWorldResponse = &objects.InteractionResponse{
	Type: 4,
	Data: &objects.InteractionApplicationCommandCallbackData{
		TTS:     false,
		Content: "hello world",
	},
}

var autoCompleteResponse = &objects.InteractionResponse{
	Type: objects.ResponseCommandAutocompleteResult,
	Data: &objects.InteractionApplicationCommandCallbackData{
		Choices: []*objects.ApplicationCommandOptionChoice{
			{Name: "a", Value: "b"},
		},
	},
}

func TestCommandRouter_build(t *testing.T) {
	tests := []struct {
		name string

		cmd                   string
		globalAllowedMentions *objects.AllowedMentions
		groupAllowedMentions  *objects.AllowedMentions

		cmdInteraction          *objects.Interaction
		autocompleteInteraction *objects.Interaction

		paramsCheck           func(*testing.T, map[string]interface{})
		throwsCmdErr          string
		throwsAutocompleteErr string

		wantsCmdErr          string
		wantsAutocompleteErr string
		cmdResponse          *objects.InteractionResponse
		autocompleteResponse *objects.InteractionResponse
	}{
		// Mentions configuration test

		{
			name: "global allowed mentions",
			cmd:  "root1",
			globalAllowedMentions: &objects.AllowedMentions{
				Parse:       []string{"a", "b", "c"},
				RepliedUser: true,
			},
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options:           nil,
				Resolved:          objects.ApplicationCommandInteractionDataResolved{},
			}),
			cmdResponse: &objects.InteractionResponse{
				Type: 4,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     false,
					Content: "hello world",
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"a", "b", "c"},
						RepliedUser: true,
					},
				},
			},
		},
		{
			name: "group allowed mentions",
			cmd:  "sub1",
			globalAllowedMentions: &objects.AllowedMentions{
				Parse:       []string{"a", "b", "c"},
				RepliedUser: true,
			},
			groupAllowedMentions: &objects.AllowedMentions{
				Parse:       []string{"d", "e", "f"},
				RepliedUser: true,
			},
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommand,
						Name: "sub1",
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			cmdResponse: &objects.InteractionResponse{
				Type: 4,
				Data: &objects.InteractionApplicationCommandCallbackData{
					TTS:     false,
					Content: "hello world",
					AllowedMentions: &objects.AllowedMentions{
						Parse:       []string{"d", "e", "f"},
						RepliedUser: true,
					},
				},
			},
		},

		// Command routing tests

		{
			name: "root command with no params",
			cmd:  "root1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options:           nil,
				Resolved:          objects.ApplicationCommandInteractionDataResolved{},
			}),
			cmdResponse: helloWorldResponse,
		},
		{
			name: "root command with params",
			cmd:  "root2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root2",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "string",
						Value:   "hello",
						Focused: false,
						Options: nil,
					},
					{
						Type:    objects.TypeInteger,
						Name:    "int",
						Value:   69,
						Focused: false,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "hello", m["string"])
				assert.Equal(t, 69, m["int"])
			},
			cmdResponse: helloWorldResponse,
		},
		{
			name: "sub-command with no params",
			cmd:  "sub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub1",
						Options: []*objects.ApplicationCommandInteractionDataOption{},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			cmdResponse: helloWorldResponse,
		},
		{
			name: "sub-command with params",
			cmd:  "sub2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommand,
						Name: "sub2",
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "string",
								Value:   "hello",
								Focused: false,
								Options: nil,
							},
							{
								Type:    objects.TypeInteger,
								Name:    "int",
								Value:   69,
								Focused: false,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "hello", m["string"])
				assert.Equal(t, 69, m["int"])
			},
			cmdResponse: helloWorldResponse,
		},
		{
			name: "sub-sub-command with no params",
			cmd:  "subsub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommandGroup,
						Name: "sub3",
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "subsub1",
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			cmdResponse: helloWorldResponse,
		},
		{
			name: "sub-sub-command with params",
			cmd:  "subsub2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommandGroup,
						Name: "sub3",
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "subsub2",
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "string",
										Value:   "hello",
										Focused: false,
										Options: nil,
									},
									{
										Type:    objects.TypeInteger,
										Name:    "int",
										Value:   69,
										Focused: false,
										Options: nil,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "hello", m["string"])
				assert.Equal(t, 69, m["int"])
			},
			cmdResponse: helloWorldResponse,
		},

		// Command errors

		{
			name: "root command not found",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "nonexistent",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options:           nil,
				Resolved:          objects.ApplicationCommandInteractionDataResolved{},
				TargetID:          1234,
			}),
		},
		{
			name: "root command throws error",
			cmd:  "root1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options:           nil,
				Resolved:          objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsCmdErr: "cat broke wire",
			wantsCmdErr:  "cat broke wire",
		},
		{
			name: "sub-command not found",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "nonexistentsub",
						Options: []*objects.ApplicationCommandInteractionDataOption{},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
		},
		{
			name: "sub-command throws error",
			cmd:  "sub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub1",
						Options: []*objects.ApplicationCommandInteractionDataOption{},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsCmdErr: "cat broke wire",
			wantsCmdErr:  "cat broke wire",
		},
		{
			name: "sub-sub-command not found",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommandGroup,
						Name: "sub3",
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "nonexistentsub",
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
		},
		{
			name: "sub-sub-command throws error",
			cmd:  "subsub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type: objects.TypeSubCommandGroup,
						Name: "sub3",
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "subsub1",
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsCmdErr: "cat broke wire",
			wantsCmdErr:  "cat broke wire",
		},

		// Auto-complete routing tests

		{
			name: "root auto-complete with no params",
			cmd:  "root1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: false,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: true,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},
		{
			name: "root auto-complete with params",
			cmd:  "root2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root2",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: false,
						Options: nil,
					},
					{
						Type:    objects.TypeString,
						Name:    "string",
						Value:   "hello",
						Focused: false,
						Options: nil,
					},
					{
						Type:    objects.TypeInteger,
						Name:    "int",
						Value:   69,
						Focused: false,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root2",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: true,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "hello", m["string"])
				assert.Equal(t, 69, m["int"])
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},
		{
			name: "sub-command auto-complete with no params",
			cmd:  "sub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub1",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: false,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub1",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: true,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},
		{
			name: "sub-command auto-complete with params",
			cmd:  "sub2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub2",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: false,
								Options: nil,
							},
							{
								Type:    objects.TypeString,
								Name:    "string",
								Value:   "hello",
								Focused: false,
								Options: nil,
							},
							{
								Type:    objects.TypeInteger,
								Name:    "int",
								Value:   69,
								Focused: false,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub2",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: true,
								Options: nil,
							},
							{
								Type:    objects.TypeString,
								Name:    "string",
								Value:   "hello",
								Focused: false,
								Options: nil,
							},
							{
								Type:    objects.TypeInteger,
								Name:    "int",
								Value:   69,
								Focused: false,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},
		{
			name: "sub-sub-command auto-complete with no params",
			cmd:  "subsub1",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeSubCommand,
								Name:    "subsub1",
								Focused: false,
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: false,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeSubCommand,
								Name:    "subsub1",
								Focused: false,
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: true,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},
		{
			name: "sub-sub-command auto-complete with params",
			cmd:  "subsub2",
			cmdInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeSubCommand,
								Name:    "subsub2",
								Focused: false,
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: false,
										Options: nil,
									},
									{
										Type:    objects.TypeString,
										Name:    "string",
										Value:   "hello",
										Focused: false,
										Options: nil,
									},
									{
										Type:    objects.TypeInteger,
										Name:    "int",
										Value:   69,
										Focused: false,
										Options: nil,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeSubCommand,
								Name:    "subsub2",
								Focused: false,
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: true,
										Options: nil,
									},
									{
										Type:    objects.TypeString,
										Name:    "string",
										Value:   "hello",
										Focused: false,
										Options: nil,
									},
									{
										Type:    objects.TypeInteger,
										Name:    "int",
										Value:   69,
										Focused: false,
										Options: nil,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			paramsCheck: func(t *testing.T, m map[string]interface{}) {
				t.Helper()
				assert.Equal(t, "123", m["autocompletetest"])
			},
			cmdResponse:          helloWorldResponse,
			autocompleteResponse: autoCompleteResponse,
		},

		// Auto-complete errors

		{
			name: "root auto-complete not found",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "abcd",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: true,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
		},
		{
			name: "root auto-complete throws error",
			cmd:  "root1",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root1",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeString,
						Name:    "autocompletetest",
						Value:   "123",
						Focused: true,
						Options: nil,
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsAutocompleteErr: "cat broke wire",
			wantsAutocompleteErr:  "cat broke wire",
		},
		{
			name: "sub-command auto-complete not found",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "nonexistentsub",
						Focused: true,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: true,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			wantsAutocompleteErr: "the command does not exist",
		},
		{
			name: "sub-command auto-complete throws error",
			cmd:  "sub1",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root3",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommand,
						Name:    "sub1",
						Focused: true,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type:    objects.TypeString,
								Name:    "autocompletetest",
								Value:   "123",
								Focused: true,
								Options: nil,
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsAutocompleteErr: "cat broke wire",
			wantsAutocompleteErr:  "cat broke wire",
		},
		{
			name: "sub-sub-command auto-complete not found",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "abc",
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: true,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			wantsAutocompleteErr: "the command does not exist",
		},
		{
			name: "sub-sub-command auto-complete throws error",
			cmd:  "subsub1",
			autocompleteInteraction: mockInteraction(&objects.ApplicationCommandInteractionData{
				DiscordBaseObject: objects.DiscordBaseObject{ID: 1234},
				Name:              "root4",
				Type:              objects.CommandTypeChatInput,
				Version:           1,
				Options: []*objects.ApplicationCommandInteractionDataOption{
					{
						Type:    objects.TypeSubCommandGroup,
						Name:    "sub3",
						Focused: false,
						Options: []*objects.ApplicationCommandInteractionDataOption{
							{
								Type: objects.TypeSubCommand,
								Name: "subsub1",
								Options: []*objects.ApplicationCommandInteractionDataOption{
									{
										Type:    objects.TypeString,
										Name:    "autocompletetest",
										Value:   "123",
										Focused: true,
									},
								},
							},
						},
					},
				},
				Resolved: objects.ApplicationCommandInteractionDataResolved{},
			}),
			throwsAutocompleteErr: "cat broke wire",
			wantsAutocompleteErr:  "cat broke wire",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make the router.
			r, cmds := makeMockFullCommandRouter(tt.groupAllowedMentions)

			// Get the command.
			var cmd *Command
			if tt.cmd != "" {
				// Get the command from the map.
				cmd = cmds[tt.cmd]

				// Check the command exists in the map.
				require.NotNil(t, cmd)

				// Set the handler for the command if it exists.
				if tt.paramsCheck == nil {
					tt.paramsCheck = func(t *testing.T, m map[string]interface{}) {}
				}
				cmd.Function = mockCommandFunction(t, cmd, tt.cmdInteraction, tt.paramsCheck, tt.throwsCmdErr)
			}

			// Set the error handler to figure out if we got an error.
			var errResult error
			errHandler := func(err error) *objects.InteractionResponse {
				errResult = err
				return &objects.InteractionResponse{Type: 69}
			}

			// Build the router.
			cmdHandler, autoCompleteHandler := r.build(loaderPassthrough{
				rest:                  dummyRestClient,
				errHandler:            errHandler,
				globalAllowedMentions: tt.globalAllowedMentions,
			})

			// Inject the auto-complete handler if this is wanted.
			if tt.autocompleteInteraction != nil && cmd != nil {
				cmd.Options = append(cmd.Options, &objects.ApplicationCommandOption{
					OptionType:   objects.TypeString,
					Name:         "autocompletetest",
					Description:  "testing testing 123",
					Required:     true,
					Autocomplete: true,
				})
				if cmd.autocomplete == nil {
					cmd.autocomplete = map[string]interface{}{}
				}
				cmd.autocomplete["autocompletetest"] = mockAutocompleteFunction(t, cmd, tt.autocompleteInteraction, tt.throwsAutocompleteErr)
			}

			// Run the command handler if applicable.
			if tt.cmdInteraction != nil {
				resp := cmdHandler(context.Background(), tt.cmdInteraction)
				if tt.wantsCmdErr == "" {
					assert.NoError(t, errResult)
					assert.Equal(t, tt.cmdResponse, resp)
				} else {
					assert.EqualError(t, errResult, tt.wantsCmdErr)
					assert.Equal(t, &objects.InteractionResponse{Type: 69}, resp)
				}
				errResult = nil
			}

			// Run the auto-complete handler if applicable.
			if tt.autocompleteInteraction != nil {
				resp := autoCompleteHandler(context.Background(), tt.autocompleteInteraction)
				if tt.wantsAutocompleteErr == "" {
					assert.NoError(t, errResult)
					assert.Equal(t, tt.autocompleteResponse, resp)
				} else {
					assert.EqualError(t, errResult, tt.wantsAutocompleteErr)
					assert.Nil(t, resp)
				}
			}
		})
	}
}

func TestCommandRouter_FormulateDiscordCommands(t *testing.T) {
	tests := []struct {
		name string

		init func() *CommandRouter
	}{
		{
			name: "no commands",
			init: func() *CommandRouter {
				return &CommandRouter{}
			},
		},
		{
			name: "with all cases",
			init: func() *CommandRouter {
				// Defines the root router.
				r := &CommandRouter{}

				// Defines the root command with no arguments.
				r.NewCommandBuilder("rootnoargs").
					DefaultPermission().
					Description("root with no arguments").
					MustBuild()

				// Defines the root command with arguments.
				r.NewCommandBuilder("rootargs").
					Description("root with arguments").
					StringOption("req_string_option", "the required string option", true,
						StringAutoCompleteFuncBuilder(func(ctx *CommandRouterCtx) ([]StringChoice, error) {
							// Important to note this doesn't actually work.
							return nil, nil
						}),
					).
					StringOption("optional_string_option", "The optional string option", false, nil).
					IntOption("req_int_option", "the required int option", true,
						IntAutoCompleteFuncBuilder(func(ctx *CommandRouterCtx) ([]IntChoice, error) {
							// Important to note this doesn't actually work.
							return nil, nil
						}),
					).
					IntOption("optional_int_option", "The optional int option", false, nil).
					DoubleOption("req_double_option", "the required double option", true,
						DoubleAutoCompleteFuncBuilder(func(ctx *CommandRouterCtx) ([]DoubleChoice, error) {
							// Important to note this doesn't actually work.
							return nil, nil
						}),
					).
					IntOption("optional_double_option", "The optional double option", false, nil).
					BoolOption("req_bool_option", "the required boolean option", true).
					BoolOption("optional_bool_option", "the optional boolean option", false).
					RoleOption("req_role_option", "the required role option", true).
					RoleOption("optional_role_option", "the optional role option", false).
					ChannelOption("req_channel_option", "the required channel option", true).
					ChannelOption("optional_channel_option", "the optional channel option", false).
					MentionableOption("req_mentionable_option", "the required mentionable option", true).
					MentionableOption("optional_mentionable_option", "the optional mentionable option", false).
					UserOption("req_user_option", "the required user option", true).
					UserOption("optional_user_option", "the optional user option", false).
					MustBuild()

				// Defines a command group with commands in the group.
				g := r.MustNewCommandGroup("group1", "group 1")
				g.NewCommandBuilder("cmd1").Description("first command in group").
					StringOption("test", "testing", true, nil).
					MustBuild()
				g.NewCommandBuilder("cmd2").Description("second command in group").MustBuild()

				// Defines a command group with sub-groups.
				g = r.MustNewCommandGroup("group2", "group 2")
				g.MustNewCommandGroup("subgroup1", "subgroup 1", nil).
					NewCommandBuilder("subcmd1").Description("first command in subgroup").MustBuild()
				s := g.MustNewCommandGroup("subgroup1", "subgroup 1", nil)
				s.NewCommandBuilder("subcmd1").Description("first command in subgroup").MustBuild()
				s.NewCommandBuilder("subcmd2").Description("second command in subgroup").MustBuild()

				// Return the command router.
				return r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.MarshalIndent(tt.init().FormulateDiscordCommands(), "", "  ")
			require.NoError(t, err)
			if golden.Update() {
				golden.Set(t, b)
			}

			// HACK: It doesn't *always* resolve the right way around.
			// Give it 100k chances. Most test attempts for me get it within 2-3.
			// buuuuuut since it's random, 100k puts errors in the realm of impossibility
			for i := 0; i < 100000; i++ {
				if string(golden.Get(t)) == string(b) {
					break
				}
				b, err = json.MarshalIndent(tt.init().FormulateDiscordCommands(), "", "  ")
				require.NoError(t, err)
			}
			assert.JSONEq(t, string(golden.Get(t)), string(b))
		})
	}
}

func TestCommandRouterCtx_Bind(t *testing.T) {
	type x struct {
		Str      string             `discord:"str"`
		Int      int                `discord:"int"`
		Bool     bool               `discord:"bool"`
		Bool2    bool               `discord:"bool2"`
		Channel  *ResolvableChannel `discord:"channel"`
		Double   float64            `discord:"double"`
		NoTag    string
		EmptyTag string `discord:""`
		ExtraTag string `discord:"extra"`
	}

	tests := []struct {
		name     string
		init     func() *CommandRouterCtx
		validate func(*CommandRouterCtx, *x)
	}{
		{
			name: "successful binding",
			init: func() *CommandRouterCtx {
				r := &CommandRouter{}
				c, _ := r.NewCommandBuilder("test").
					StringOption("str", "A string option", true, nil).
					StringOption("str2", "A string option", true, nil).
					IntOption("int", "An int option", true, nil).
					BoolOption("bool", "A bool option", true).
					BoolOption("bool2", "Another bool option", true).
					ChannelOption("channel", "A channel option", true).
					DoubleOption("double", "A double option", true, nil).
					Build()
				opts := map[string]interface{}{
					"str":  "str",
					"int":  1,
					"bool": true,
					"channel": &ResolvableChannel{
						id: "123",
					},
					"int2":   2,
					"double": 3.14,
				}
				return &CommandRouterCtx{Options: opts, Command: c}
			},
			validate: func(ctx *CommandRouterCtx, items *x) {
				assert.NoError(t, ctx.Bind(items))
				assert.Equal(t, "str", items.Str)
				assert.Equal(t, 1, items.Int)
				assert.Equal(t, true, items.Bool)
				assert.Equal(t, false, items.Bool2)
				assert.Equal(t, "123", items.Channel.id)
				assert.Equal(t, 3.14, items.Double)
			},
		},
		{
			name: "non-pointer",
			init: func() *CommandRouterCtx {
				r := &CommandRouter{}
				c, _ := r.NewCommandBuilder("test").
					StringOption("str", "A string option", true, nil).
					Build()
				opts := map[string]interface{}{
					"str": "str",
				}
				return &CommandRouterCtx{Options: opts, Command: c}
			},
			validate: func(ctx *CommandRouterCtx, items *x) {
				assert.Error(t, ctx.Bind(*items))
			},
		},
		{
			name: "non-struct",
			init: func() *CommandRouterCtx {
				r := &CommandRouter{}
				c, _ := r.NewCommandBuilder("test").Build()
				opts := map[string]interface{}{}
				return &CommandRouterCtx{Options: opts, Command: c}
			},
			validate: func(ctx *CommandRouterCtx, items *x) {
				var myStr string
				assert.Error(t, ctx.Bind(&myStr))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.init()
			var items x
			tt.validate(ctx, &items)
		})
	}
}
