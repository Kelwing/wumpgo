// Code generated by generate_command_builder.go; DO NOT EDIT.

package router

//go:generate go run generate_command_builder.go

import (
	"fmt"

	"wumpgo.dev/wumpgo/objects"
)

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

func (c *commandBuilder[T]) StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) T {
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
			c.cmd.autocomplete = map[string]any{}
		}
		c.cmd.autocomplete[name] = f
	}
	return builderWrapify(c)
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

func (c *commandBuilder[T]) IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) T {
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
			c.cmd.autocomplete = map[string]any{}
		}
		c.cmd.autocomplete[name] = f
	}
	return builderWrapify(c)
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

func (c *commandBuilder[T]) DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) T {
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
			c.cmd.autocomplete = map[string]any{}
		}
		c.cmd.autocomplete[name] = f
	}
	return builderWrapify(c)
}

type textCommandBuilder struct {
	*commandBuilder[TextCommandBuilder]
}

func (c *commandBuilder[T]) TextCommand() TextCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeChatInput)
	return textCommandBuilder{(*commandBuilder[TextCommandBuilder])(c)}
}

type subcommandBuilder struct {
	*commandBuilder[SubCommandBuilder]
}

type messageCommandBuilder struct {
	*commandBuilder[MessageCommandBuilder]
}

func (c *commandBuilder[T]) MessageCommand() MessageCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeMessage)
	return messageCommandBuilder{(*commandBuilder[MessageCommandBuilder])(c)}
}

type userCommandBuilder struct {
	*commandBuilder[UserCommandBuilder]
}

func (c *commandBuilder[T]) UserCommand() UserCommandBuilder {
	c.cmd.commandType = int(objects.CommandTypeUser)
	return userCommandBuilder{(*commandBuilder[UserCommandBuilder])(c)}
}

func builderWrapify[T any](c *commandBuilder[T]) T {
	var ptr *T
	switch (any)(ptr).(type) {
		case *CommandBuilder:
			return (any)(c).(T)
		case *TextCommandBuilder:
			return (any)(textCommandBuilder{(any)(c).(*commandBuilder[TextCommandBuilder])}).(T)
		case *SubCommandBuilder:
			return (any)(subcommandBuilder{(any)(c).(*commandBuilder[SubCommandBuilder])}).(T)
		case *MessageCommandBuilder:
			return (any)(messageCommandBuilder{(any)(c).(*commandBuilder[MessageCommandBuilder])}).(T)
		case *UserCommandBuilder:
			return (any)(userCommandBuilder{(any)(c).(*commandBuilder[UserCommandBuilder])}).(T)
		default:
			panic(fmt.Errorf("unknown handler: %T", ptr))
	}
}
