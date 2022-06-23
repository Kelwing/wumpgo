package router

import (
	"testing"

	"github.com/Postcord/objects"
	"github.com/stretchr/testify/assert"
)

func TestStringStaticChoicesBuilder(t *testing.T) {
	items := []StringChoice{{Name: "a"}}
	good := func(s []StringChoice) {
		t.Helper()
		assert.Equal(t, items, s)
	}
	bad := func(_ StringAutoCompleteFunc) {
		t.Fatal("should not be called")
	}
	builder := StringStaticChoicesBuilder(items)
	builder(good, bad)
}

func TestStringAutoCompleteFuncBuilder(t *testing.T) {
	f := func(_ *CommandRouterCtx) ([]StringChoice, error) {
		return nil, nil
	}
	bad := func(s []StringChoice) {
		t.Fatal("should not be called")
	}
	good := func(funcPassed StringAutoCompleteFunc) {
		t.Helper()
		assert.NotNil(t, funcPassed)
	}
	builder := StringAutoCompleteFuncBuilder(f)
	builder(bad, good)
}

func Test_commandBuilder_StringOption(t *testing.T) {
	tests := []struct {
		name string

		choices bool
		f       bool

		wantErr string
	}{
		{
			name:    "choices",
			choices: true,
			f:       false,
		},
		{
			name:    "auto-complete",
			choices: false,
			f:       true,
		},
		{
			name:    "choices and auto-complete",
			choices: true,
			f:       true,
			wantErr: "cannot set both function and choice slice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b CommandBuilder = &commandBuilder[CommandBuilder]{}
			var f StringChoiceBuilder
			if tt.choices || tt.f {
				f = func(addStaticOptions func([]StringChoice), addAutocomplete func(StringAutoCompleteFunc)) {
					if tt.choices {
						addStaticOptions([]StringChoice{
							{Name: "a", Value: "1"},
							{Name: "b", Value: "2"},
							{Name: "c", Value: "3"},
						})
					}
					if tt.f {
						addAutocomplete(func(ctx *CommandRouterCtx) ([]StringChoice, error) {
							return nil, nil
						})
					}
				}
			}
			err := callBuilderFunction(t, b, "StringOption", "testing", "testing 123", true, f)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				var discordifiedChoices []objects.ApplicationCommandOptionChoice
				if tt.choices {
					discordifiedChoices = []objects.ApplicationCommandOptionChoice{
						{Name: "a", Value: "1"},
						{Name: "b", Value: "2"},
						{Name: "c", Value: "3"},
					}
				}
				assert.Equal(t, []*objects.ApplicationCommandOption{
					{
						OptionType:   objects.TypeString,
						Name:         "testing",
						Description:  "testing 123",
						Required:     true,
						Choices:      discordifiedChoices,
						Autocomplete: tt.f,
					},
				}, b.(*commandBuilder[CommandBuilder]).cmd.Options)
				if tt.f {
					assert.IsType(t, (StringAutoCompleteFunc)(nil),
						b.(*commandBuilder[CommandBuilder]).cmd.autocomplete["testing"])
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestIntStaticChoicesBuilder(t *testing.T) {
	items := []IntChoice{{Name: "a"}}
	good := func(s []IntChoice) {
		t.Helper()
		assert.Equal(t, items, s)
	}
	bad := func(_ IntAutoCompleteFunc) {
		t.Fatal("should not be called")
	}
	builder := IntStaticChoicesBuilder(items)
	builder(good, bad)
}

func TestIntAutoCompleteFuncBuilder(t *testing.T) {
	f := func(_ *CommandRouterCtx) ([]IntChoice, error) {
		return nil, nil
	}
	bad := func(s []IntChoice) {
		t.Fatal("should not be called")
	}
	good := func(funcPassed IntAutoCompleteFunc) {
		t.Helper()
		assert.NotNil(t, funcPassed)
	}
	builder := IntAutoCompleteFuncBuilder(f)
	builder(bad, good)
}

func Test_commandBuilder_IntOption(t *testing.T) {
	tests := []struct {
		name string

		choices bool
		f       bool

		wantErr string
	}{
		{
			name:    "choices",
			choices: true,
			f:       false,
		},
		{
			name:    "auto-complete",
			choices: false,
			f:       true,
		},
		{
			name:    "choices and auto-complete",
			choices: true,
			f:       true,
			wantErr: "cannot set both function and choice slice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b CommandBuilder = &commandBuilder[CommandBuilder]{}
			var f IntChoiceBuilder
			if tt.choices || tt.f {
				f = func(addStaticOptions func([]IntChoice), addAutocomplete func(IntAutoCompleteFunc)) {
					if tt.choices {
						addStaticOptions([]IntChoice{
							{Name: "a", Value: 1},
							{Name: "b", Value: 2},
							{Name: "c", Value: 3},
						})
					}
					if tt.f {
						addAutocomplete(func(ctx *CommandRouterCtx) ([]IntChoice, error) {
							return nil, nil
						})
					}
				}
			}
			err := callBuilderFunction(t, b, "IntOption", "testing", "testing 123", true, f)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				var discordifiedChoices []objects.ApplicationCommandOptionChoice
				if tt.choices {
					discordifiedChoices = []objects.ApplicationCommandOptionChoice{
						{Name: "a", Value: 1},
						{Name: "b", Value: 2},
						{Name: "c", Value: 3},
					}
				}
				assert.Equal(t, []*objects.ApplicationCommandOption{
					{
						OptionType:   objects.TypeInteger,
						Name:         "testing",
						Description:  "testing 123",
						Required:     true,
						Choices:      discordifiedChoices,
						Autocomplete: tt.f,
					},
				}, b.(*commandBuilder[CommandBuilder]).cmd.Options)
				if tt.f {
					assert.IsType(t, (IntAutoCompleteFunc)(nil),
						b.(*commandBuilder[CommandBuilder]).cmd.autocomplete["testing"])
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func TestDoubleStaticChoicesBuilder(t *testing.T) {
	items := []DoubleChoice{{Name: "a"}}
	good := func(s []DoubleChoice) {
		t.Helper()
		assert.Equal(t, items, s)
	}
	bad := func(_ DoubleAutoCompleteFunc) {
		t.Fatal("should not be called")
	}
	builder := DoubleStaticChoicesBuilder(items)
	builder(good, bad)
}

func TestDoubleAutoCompleteFuncBuilder(t *testing.T) {
	f := func(_ *CommandRouterCtx) ([]DoubleChoice, error) {
		return nil, nil
	}
	bad := func(s []DoubleChoice) {
		t.Fatal("should not be called")
	}
	good := func(funcPassed DoubleAutoCompleteFunc) {
		t.Helper()
		assert.NotNil(t, funcPassed)
	}
	builder := DoubleAutoCompleteFuncBuilder(f)
	builder(bad, good)
}

func Test_commandBuilder_DoubleOption(t *testing.T) {
	tests := []struct {
		name string

		choices bool
		f       bool

		wantErr string
	}{
		{
			name:    "choices",
			choices: true,
			f:       false,
		},
		{
			name:    "auto-complete",
			choices: false,
			f:       true,
		},
		{
			name:    "choices and auto-complete",
			choices: true,
			f:       true,
			wantErr: "cannot set both function and choice slice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b CommandBuilder = &commandBuilder[CommandBuilder]{}
			var f DoubleChoiceBuilder
			if tt.choices || tt.f {
				f = func(addStaticOptions func([]DoubleChoice), addAutocomplete func(DoubleAutoCompleteFunc)) {
					if tt.choices {
						addStaticOptions([]DoubleChoice{
							{Name: "a", Value: 1},
							{Name: "b", Value: 2},
							{Name: "c", Value: 3},
						})
					}
					if tt.f {
						addAutocomplete(func(ctx *CommandRouterCtx) ([]DoubleChoice, error) {
							return nil, nil
						})
					}
				}
			}
			err := callBuilderFunction(t, b, "DoubleOption", "testing", "testing 123", true, f)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				var discordifiedChoices []objects.ApplicationCommandOptionChoice
				if tt.choices {
					discordifiedChoices = []objects.ApplicationCommandOptionChoice{
						{Name: "a", Value: (float64)(1)},
						{Name: "b", Value: (float64)(2)},
						{Name: "c", Value: (float64)(3)},
					}
				}
				assert.Equal(t, []*objects.ApplicationCommandOption{
					{
						OptionType:   objects.TypeNumber,
						Name:         "testing",
						Description:  "testing 123",
						Required:     true,
						Choices:      discordifiedChoices,
						Autocomplete: tt.f,
					},
				}, b.(*commandBuilder[CommandBuilder]).cmd.Options)
				if tt.f {
					assert.IsType(t, (DoubleAutoCompleteFunc)(nil),
						b.(*commandBuilder[CommandBuilder]).cmd.autocomplete["testing"])
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func Test_textCommandBuilder_StringOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "StringOption", "testing", "testing 123", true,
		StringStaticChoicesBuilder([]StringChoice{
			{Name: "a", Value: "b"},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeString,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: "b"},
			},
			Required: true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_IntOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "IntOption", "testing", "testing 123", true,
		IntStaticChoicesBuilder([]IntChoice{
			{Name: "a", Value: 123},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeInteger,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: 123},
			},
			Required: true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_BoolOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "BoolOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeBoolean,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_UserOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "UserOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeUser,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_ChannelOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "ChannelOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeChannel,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_RoleOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "RoleOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeRole,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_MentionableOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "MentionableOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeMentionable,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_DoubleOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "DoubleOption", "testing", "testing 123", true,
		DoubleStaticChoicesBuilder([]DoubleChoice{
			{Name: "a", Value: (float64)(123)},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeNumber,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: (float64)(123)},
			},
			Required: true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_AttachmentOption(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "AttachmentOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeAttachment,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(textCommandBuilder).cmd.Options)
}

func Test_textCommandBuilder_DefaultPermission(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "DefaultPermission"))
	assert.True(t, b.(textCommandBuilder).cmd.DefaultPermission)
}

func Test_textCommandBuilder_AllowedMentions(t *testing.T) {
	x := &objects.AllowedMentions{}
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "AllowedMentions", x))
	assert.Same(t, x, b.(textCommandBuilder).cmd.AllowedMentions)
}

func Test_commandBuilder_TextCommand(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	res := base.TextCommand()
	assert.NotNil(t, res.(textCommandBuilder).commandBuilder)
}

func Test_subcommandBuilder_StringOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "StringOption", "testing", "testing 123", true,
		StringStaticChoicesBuilder([]StringChoice{
			{Name: "a", Value: "b"},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeString,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: "b"},
			},
			Required: true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_IntOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "IntOption", "testing", "testing 123", true,
		IntStaticChoicesBuilder([]IntChoice{
			{Name: "a", Value: 123},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeInteger,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: 123},
			},
			Required: true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_BoolOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "BoolOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeBoolean,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_UserOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "UserOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeUser,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_ChannelOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "ChannelOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeChannel,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_RoleOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "RoleOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeRole,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_MentionableOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "MentionableOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeMentionable,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_DoubleOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "DoubleOption", "testing", "testing 123", true,
		DoubleStaticChoicesBuilder([]DoubleChoice{
			{Name: "a", Value: (float64)(123)},
		}),
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeNumber,
			Name:        "testing",
			Description: "testing 123",
			Choices: []objects.ApplicationCommandOptionChoice{
				{Name: "a", Value: (float64)(123)},
			},
			Required: true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_AttachmentOption(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(
		t, b, "AttachmentOption", "testing", "testing 123", true,
	))
	assert.Equal(t, []*objects.ApplicationCommandOption{
		{
			OptionType:  objects.TypeAttachment,
			Name:        "testing",
			Description: "testing 123",
			Required:    true,
		},
	}, b.(subcommandBuilder).cmd.Options)
}

func Test_subcommandBuilder_DefaultPermission(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "DefaultPermission"))
	assert.True(t, b.(subcommandBuilder).cmd.DefaultPermission)
}

func Test_subcommandBuilder_AllowedMentions(t *testing.T) {
	x := &objects.AllowedMentions{}
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "AllowedMentions", x))
	assert.Same(t, x, b.(subcommandBuilder).cmd.AllowedMentions)
}

func Test_messageCommandBuilder_DefaultPermission(t *testing.T) {
	var b MessageCommandBuilder = messageCommandBuilder{&commandBuilder[MessageCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "DefaultPermission"))
	assert.True(t, b.(messageCommandBuilder).cmd.DefaultPermission)
}

func Test_messageCommandBuilder_AllowedMentions(t *testing.T) {
	x := &objects.AllowedMentions{}
	var b MessageCommandBuilder = messageCommandBuilder{&commandBuilder[MessageCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "AllowedMentions", x))
	assert.Same(t, x, b.(messageCommandBuilder).cmd.AllowedMentions)
}

func Test_commandBuilder_MessageCommand(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	res := base.MessageCommand()
	assert.NotNil(t, res.(messageCommandBuilder).commandBuilder)
}

func Test_userCommandBuilder_DefaultPermission(t *testing.T) {
	var b UserCommandBuilder = userCommandBuilder{&commandBuilder[UserCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "DefaultPermission"))
	assert.True(t, b.(userCommandBuilder).cmd.DefaultPermission)
}

func Test_userCommandBuilder_AllowedMentions(t *testing.T) {
	x := &objects.AllowedMentions{}
	var b UserCommandBuilder = userCommandBuilder{&commandBuilder[UserCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "AllowedMentions", x))
	assert.Same(t, x, b.(userCommandBuilder).cmd.AllowedMentions)
}

func Test_userCommandBuilder_MessageCommand(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	res := base.UserCommand()
	assert.NotNil(t, res.(userCommandBuilder).commandBuilder)
}

func Test_commandBuilder_Description(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	assert.NoError(t, callBuilderFunction(t, base, "Description", "testing"))
	assert.Equal(t, "testing", base.cmd.Description)
}

func Test_commandBuilder_DefaultPermission(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	assert.NoError(t, callBuilderFunction(t, base, "DefaultPermission"))
	assert.True(t, base.cmd.DefaultPermission)
}

func Test_commandBuilder_AllowedMentions(t *testing.T) {
	x := &objects.AllowedMentions{}
	base := &commandBuilder[CommandBuilder]{}
	assert.NoError(t, callBuilderFunction(t, base, "AllowedMentions", x))
	assert.Same(t, x, base.cmd.AllowedMentions)
}

func Test_commandBuilder_Handler(t *testing.T) {
	base := &commandBuilder[CommandBuilder]{}
	assert.NoError(t, callBuilderFunction(t, base, "Handler", func(ctx *CommandRouterCtx) error {
		return nil
	}))
	assert.NotNil(t, base.cmd.Function)
}

func Test_textCommandBuilder_Description(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Description", "testing"))
	assert.Equal(t, "testing", b.(textCommandBuilder).cmd.Description)
}

func Test_textCommandBuilder_Handler(t *testing.T) {
	var b TextCommandBuilder = textCommandBuilder{&commandBuilder[TextCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Handler", func(ctx *CommandRouterCtx) error {
		return nil
	}))
	assert.NotNil(t, b.(textCommandBuilder).cmd.Function)
}

func Test_subcommandBuilder_Description(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Description", "testing"))
	assert.Equal(t, "testing", b.(subcommandBuilder).cmd.Description)
}

func Test_subcommandBuilder_Handler(t *testing.T) {
	var b SubCommandBuilder = subcommandBuilder{&commandBuilder[SubCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Handler", func(ctx *CommandRouterCtx) error {
		return nil
	}))
	assert.NotNil(t, b.(subcommandBuilder).cmd.Function)
}

func Test_messageCommandBuilder_Handler(t *testing.T) {
	var b MessageCommandBuilder = messageCommandBuilder{&commandBuilder[MessageCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Handler", func(ctx *CommandRouterCtx, msg *objects.Message) error {
		return nil
	}))
	assert.NotNil(t, b.(messageCommandBuilder).cmd.Function)
}

func Test_userCommandBuilder_Handler(t *testing.T) {
	var b UserCommandBuilder = userCommandBuilder{&commandBuilder[UserCommandBuilder]{}}
	assert.NoError(t, callBuilderFunction(t, b, "Handler", func(ctx *CommandRouterCtx, user *objects.GuildMember) error {
		return nil
	}))
	assert.NotNil(t, b.(userCommandBuilder).cmd.Function)
}
