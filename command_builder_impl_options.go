package router

import "github.com/Postcord/objects"

// StringChoice is used to define a choice of the string type.
type StringChoice struct {
	// Name is the name of the choice.
	Name string `json:"name"`

	// Value is the string that is the resulting value.
	Value string `json:"value"`
}

func (c *commandBuilder) StringOption(name, description string, required bool, choices []StringChoice) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	if choices != nil && len(choices) != 0 {
		discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
		for i, v := range choices {
			discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
		}
	}
	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:  objects.TypeString,
		Name:        name,
		Description: description,
		Required:    required,
		Choices:     discordifiedChoices,
	})
	return c
}

// IntChoice is used to define a choice of the int type.
type IntChoice struct {
	// Name is the name of the choice.
	Name string `json:"name"`

	// Value is the int that is the resulting value.
	Value int `json:"value"`
}

func (c *commandBuilder) IntOption(name, description string, required bool, choices []IntChoice) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	if choices != nil && len(choices) != 0 {
		discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
		for i, v := range choices {
			discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
		}
	}
	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:  objects.TypeInteger,
		Name:        name,
		Description: description,
		Required:    required,
		Choices:     discordifiedChoices,
	})
	return c
}

func (c *commandBuilder) appendOption(type_ objects.ApplicationCommandOptionType, name, description string, required, default_ bool) CommandBuilder {
	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:  type_,
		Name:        name,
		Description: description,
		Required:    required,
		Default:     default_,
	})
	return c
}

func (c *commandBuilder) BoolOption(name, description string, required, default_ bool) CommandBuilder {
	return c.appendOption(objects.TypeBoolean, name, description, required, default_)
}

func (c *commandBuilder) UserOption(name, description string, required bool) CommandBuilder {
	return c.appendOption(objects.TypeUser, name, description, required, false)
}

func (c *commandBuilder) ChannelOption(name, description string, required bool) CommandBuilder {
	return c.appendOption(objects.TypeChannel, name, description, required, false)
}

func (c *commandBuilder) RoleOption(name, description string, required bool) CommandBuilder {
	return c.appendOption(objects.TypeRole, name, description, required, false)
}

func (c *commandBuilder) MentionableOption(name, description string, required bool) CommandBuilder {
	return c.appendOption(objects.TypeMentionable, name, description, required, false)
}

// DoubleChoice is used to define a choice of the double type.
type DoubleChoice struct {
	// Name is the name of the choice.
	Name string `json:"name"`

	// Value is the double that is the resulting value.
	Value float64 `json:"value"`
}

func (c *commandBuilder) DoubleOption(name, description string, required bool, choices []DoubleChoice) CommandBuilder {
	var discordifiedChoices []objects.ApplicationCommandOptionChoice
	if choices != nil && len(choices) != 0 {
		discordifiedChoices = make([]objects.ApplicationCommandOptionChoice, len(choices))
		for i, v := range choices {
			discordifiedChoices[i] = objects.ApplicationCommandOptionChoice{Name: v.Name, Value: v.Value}
		}
	}
	c.cmd.Options = append(c.cmd.Options, &objects.ApplicationCommandOption{
		OptionType:  objects.TypeDouble,
		Name:        name,
		Description: description,
		Required:    required,
		Choices:     discordifiedChoices,
	})
	return c
}
