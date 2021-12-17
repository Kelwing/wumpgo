// Code generated by generate_command_builder.go; DO NOT EDIT.

package router

//go:generate go run generate_command_builder.go

import "github.com/Postcord/objects"

// StringAutoCompleteFunc is used to define the auto-complete function for StringChoice.
// Note that the command context is a special case in that the response is not used.
type StringAutoCompleteFunc = func(*CommandRouterCtx) ([]StringChoice, error)

// StringChoiceBuilder is used to choose how this choice is handled.
// This can be nil, or it can pass to one of the functions. The first function adds static choices to the router. The second option adds an autocomplete function.
// Note that you cannot call both functions.
type StringChoiceBuilder = func(addStaticOptions func([]StringChoice), addAutocomplete func(StringAutoCompleteFunc))

// StringStaticChoicesBuilder is used to create a shorthand for adding choices.
func StringStaticChoicesBuilder(choices []StringChoice) StringChoiceBuilder {
	return func(addStaticOptions func([]StringChoice), _ func(StringAutoCompleteFunc)) {
		addStaticOptions(choices)
	}
}

// StringAutoCompleteFuncBuilder is used to create a shorthand for adding a auto-complete function.
func StringAutoCompleteFuncBuilder(f StringAutoCompleteFunc) StringChoiceBuilder {
	return func(_ func([]StringChoice), addAutocomplete func(StringAutoCompleteFunc)) {
		addAutocomplete(f)
	}
}

func (c *commandBuilder) StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	var f StringAutoCompleteFunc
	if choiceBuilder != nil {
		choiceBuilder(func(choices []StringChoice) {
			if f != nil {
				panic("cannot set both function and choice slice")
			}
			discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
			for i, v := range choices {
				discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
			}
		}, func(autoCompleteFunc StringAutoCompleteFunc) {
			if discordifiedChoices != nil {
				panic("cannot set both function and choice slice")
			}
			f = autoCompleteFunc
		})
	}

	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:   objects.TypeString,
		Name:         name,
		Description:  description,
		Required:     required,
		Choices:      discordifiedChoices,
		Autocomplete: f != nil,
	})
	if f != nil {
		if c.cmd.autocomplete == nil {
			c.cmd.autocomplete = map[string]interface{}{}
		}
		c.cmd.autocomplete[name] = f
	}
	return c
}

// IntAutoCompleteFunc is used to define the auto-complete function for IntChoice.
// Note that the command context is a special case in that the response is not used.
type IntAutoCompleteFunc = func(*CommandRouterCtx) ([]IntChoice, error)

// IntChoiceBuilder is used to choose how this choice is handled.
// This can be nil, or it can pass to one of the functions. The first function adds static choices to the router. The second option adds an autocomplete function.
// Note that you cannot call both functions.
type IntChoiceBuilder = func(addStaticOptions func([]IntChoice), addAutocomplete func(IntAutoCompleteFunc))

// IntStaticChoicesBuilder is used to create a shorthand for adding choices.
func IntStaticChoicesBuilder(choices []IntChoice) IntChoiceBuilder {
	return func(addStaticOptions func([]IntChoice), _ func(IntAutoCompleteFunc)) {
		addStaticOptions(choices)
	}
}

// IntAutoCompleteFuncBuilder is used to create a shorthand for adding a auto-complete function.
func IntAutoCompleteFuncBuilder(f IntAutoCompleteFunc) IntChoiceBuilder {
	return func(_ func([]IntChoice), addAutocomplete func(IntAutoCompleteFunc)) {
		addAutocomplete(f)
	}
}

func (c *commandBuilder) IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	var f IntAutoCompleteFunc
	if choiceBuilder != nil {
		choiceBuilder(func(choices []IntChoice) {
			if f != nil {
				panic("cannot set both function and choice slice")
			}
			discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
			for i, v := range choices {
				discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
			}
		}, func(autoCompleteFunc IntAutoCompleteFunc) {
			if discordifiedChoices != nil {
				panic("cannot set both function and choice slice")
			}
			f = autoCompleteFunc
		})
	}

	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:   objects.TypeInteger,
		Name:         name,
		Description:  description,
		Required:     required,
		Choices:      discordifiedChoices,
		Autocomplete: f != nil,
	})
	if f != nil {
		if c.cmd.autocomplete == nil {
			c.cmd.autocomplete = map[string]interface{}{}
		}
		c.cmd.autocomplete[name] = f
	}
	return c
}

// DoubleAutoCompleteFunc is used to define the auto-complete function for DoubleChoice.
// Note that the command context is a special case in that the response is not used.
type DoubleAutoCompleteFunc = func(*CommandRouterCtx) ([]DoubleChoice, error)

// DoubleChoiceBuilder is used to choose how this choice is handled.
// This can be nil, or it can pass to one of the functions. The first function adds static choices to the router. The second option adds an autocomplete function.
// Note that you cannot call both functions.
type DoubleChoiceBuilder = func(addStaticOptions func([]DoubleChoice), addAutocomplete func(DoubleAutoCompleteFunc))

// DoubleStaticChoicesBuilder is used to create a shorthand for adding choices.
func DoubleStaticChoicesBuilder(choices []DoubleChoice) DoubleChoiceBuilder {
	return func(addStaticOptions func([]DoubleChoice), _ func(DoubleAutoCompleteFunc)) {
		addStaticOptions(choices)
	}
}

// DoubleAutoCompleteFuncBuilder is used to create a shorthand for adding a auto-complete function.
func DoubleAutoCompleteFuncBuilder(f DoubleAutoCompleteFunc) DoubleChoiceBuilder {
	return func(_ func([]DoubleChoice), addAutocomplete func(DoubleAutoCompleteFunc)) {
		addAutocomplete(f)
	}
}

func (c *commandBuilder) DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	var f DoubleAutoCompleteFunc
	if choiceBuilder != nil {
		choiceBuilder(func(choices []DoubleChoice) {
			if f != nil {
				panic("cannot set both function and choice slice")
			}
			discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
			for i, v := range choices {
				discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
			}
		}, func(autoCompleteFunc DoubleAutoCompleteFunc) {
			if discordifiedChoices != nil {
				panic("cannot set both function and choice slice")
			}
			f = autoCompleteFunc
		})
	}

	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:   objects.TypeNumber,
		Name:         name,
		Description:  description,
		Required:     required,
		Choices:      discordifiedChoices,
		Autocomplete: f != nil,
	})
	if f != nil {
		if c.cmd.autocomplete == nil {
			c.cmd.autocomplete = map[string]interface{}{}
		}
		c.cmd.autocomplete[name] = f
	}
	return c
}

type textCommandBuilder struct {
	*commandBuilder
}

func (c textCommandBuilder) StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) TextCommandBuilder {
	c.commandBuilder.StringOption(name, description, required, choiceBuilder)
	return c
}

func (c textCommandBuilder) IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) TextCommandBuilder {
	c.commandBuilder.IntOption(name, description, required, choiceBuilder)
	return c
}

func (c textCommandBuilder) BoolOption(name, description string, required bool) TextCommandBuilder {
	c.commandBuilder.BoolOption(name, description, required)
	return c
}

func (c textCommandBuilder) UserOption(name, description string, required bool) TextCommandBuilder {
	c.commandBuilder.UserOption(name, description, required)
	return c
}

func (c textCommandBuilder) ChannelOption(name, description string, required bool) TextCommandBuilder {
	c.commandBuilder.ChannelOption(name, description, required)
	return c
}

func (c textCommandBuilder) RoleOption(name, description string, required bool) TextCommandBuilder {
	c.commandBuilder.RoleOption(name, description, required)
	return c
}

func (c textCommandBuilder) MentionableOption(name, description string, required bool) TextCommandBuilder {
	c.commandBuilder.MentionableOption(name, description, required)
	return c
}

func (c textCommandBuilder) DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) TextCommandBuilder {
	c.commandBuilder.DoubleOption(name, description, required, choiceBuilder)
	return c
}

func (c textCommandBuilder) DefaultPermission() TextCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c textCommandBuilder) AllowedMentions(config *objects.AllowedMentions) TextCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) TextCommand() TextCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeChatInput)
	return textCommandBuilder{c}
}

type subcommandBuilder struct {
	*commandBuilder
}

func (c subcommandBuilder) StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) SubCommandBuilder {
	c.commandBuilder.StringOption(name, description, required, choiceBuilder)
	return c
}

func (c subcommandBuilder) IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) SubCommandBuilder {
	c.commandBuilder.IntOption(name, description, required, choiceBuilder)
	return c
}

func (c subcommandBuilder) BoolOption(name, description string, required bool) SubCommandBuilder {
	c.commandBuilder.BoolOption(name, description, required)
	return c
}

func (c subcommandBuilder) UserOption(name, description string, required bool) SubCommandBuilder {
	c.commandBuilder.UserOption(name, description, required)
	return c
}

func (c subcommandBuilder) ChannelOption(name, description string, required bool) SubCommandBuilder {
	c.commandBuilder.ChannelOption(name, description, required)
	return c
}

func (c subcommandBuilder) RoleOption(name, description string, required bool) SubCommandBuilder {
	c.commandBuilder.RoleOption(name, description, required)
	return c
}

func (c subcommandBuilder) MentionableOption(name, description string, required bool) SubCommandBuilder {
	c.commandBuilder.MentionableOption(name, description, required)
	return c
}

func (c subcommandBuilder) DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) SubCommandBuilder {
	c.commandBuilder.DoubleOption(name, description, required, choiceBuilder)
	return c
}

func (c subcommandBuilder) DefaultPermission() SubCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c subcommandBuilder) AllowedMentions(config *objects.AllowedMentions) SubCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

type messageCommandBuilder struct {
	*commandBuilder
}

func (c messageCommandBuilder) DefaultPermission() MessageCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c messageCommandBuilder) AllowedMentions(config *objects.AllowedMentions) MessageCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) MessageCommand() MessageCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeMessage)
	return messageCommandBuilder{c}
}

type userCommandBuilder struct {
	*commandBuilder
}

func (c userCommandBuilder) DefaultPermission() UserCommandBuilder {
	c.commandBuilder.DefaultPermission()
	return c
}

func (c userCommandBuilder) AllowedMentions(config *objects.AllowedMentions) UserCommandBuilder {
	c.commandBuilder.AllowedMentions(config)
	return c
}

func (c *commandBuilder) UserCommand() UserCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeUser)
	return userCommandBuilder{c}
}

type commandOptions interface {
	// StringOption is used to define an option of the type string. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 3 (STRING): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) CommandBuilder

	// IntOption is used to define an option of the type int. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 4 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) CommandBuilder

	// IntOption is used to define an option of the type bool.
	// Maps to option type 5 (BOOLEAN): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	BoolOption(name, description string, required bool) CommandBuilder

	// IntOption is used to define an option of the type user.
	// Maps to option type 6 (USER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	UserOption(name, description string, required bool) CommandBuilder

	// ChannelOption is used to define an option of the type channel.
	// Maps to option type 7 (CHANNEL): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	ChannelOption(name, description string, required bool) CommandBuilder

	// RoleOption is used to define an option of the type role.
	// Maps to option type 8 (ROLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	RoleOption(name, description string, required bool) CommandBuilder

	// MentionableOption is used to define an option of the type mentionable.
	// Maps to option type 9 (MENTIONABLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	MentionableOption(name, description string, required bool) CommandBuilder

	// DoubleOption is used to define an option of the type double. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 10 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) CommandBuilder
}

type subCommandOptions interface {
	// StringOption is used to define an option of the type string. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 3 (STRING): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) SubCommandBuilder

	// IntOption is used to define an option of the type int. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 4 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) SubCommandBuilder

	// IntOption is used to define an option of the type bool.
	// Maps to option type 5 (BOOLEAN): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	BoolOption(name, description string, required bool) SubCommandBuilder

	// IntOption is used to define an option of the type user.
	// Maps to option type 6 (USER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	UserOption(name, description string, required bool) SubCommandBuilder

	// ChannelOption is used to define an option of the type channel.
	// Maps to option type 7 (CHANNEL): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	ChannelOption(name, description string, required bool) SubCommandBuilder

	// RoleOption is used to define an option of the type role.
	// Maps to option type 8 (ROLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	RoleOption(name, description string, required bool) SubCommandBuilder

	// MentionableOption is used to define an option of the type mentionable.
	// Maps to option type 9 (MENTIONABLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	MentionableOption(name, description string, required bool) SubCommandBuilder

	// DoubleOption is used to define an option of the type double. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 10 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) SubCommandBuilder
}

type textCommandOptions interface {
	// StringOption is used to define an option of the type string. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 3 (STRING): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) TextCommandBuilder

	// IntOption is used to define an option of the type int. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 4 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) TextCommandBuilder

	// IntOption is used to define an option of the type bool.
	// Maps to option type 5 (BOOLEAN): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	BoolOption(name, description string, required bool) TextCommandBuilder

	// IntOption is used to define an option of the type user.
	// Maps to option type 6 (USER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	UserOption(name, description string, required bool) TextCommandBuilder

	// ChannelOption is used to define an option of the type channel.
	// Maps to option type 7 (CHANNEL): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	ChannelOption(name, description string, required bool) TextCommandBuilder

	// RoleOption is used to define an option of the type role.
	// Maps to option type 8 (ROLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	RoleOption(name, description string, required bool) TextCommandBuilder

	// MentionableOption is used to define an option of the type mentionable.
	// Maps to option type 9 (MENTIONABLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	MentionableOption(name, description string, required bool) TextCommandBuilder

	// DoubleOption is used to define an option of the type double. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 10 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) TextCommandBuilder
}
