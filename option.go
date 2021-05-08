package interactions

import (
	"strconv"

	"github.com/Postcord/objects"
)

type CommandOption struct {
	optionType objects.ApplicationCommandOptionType
	value      interface{}
	options    []*objects.ApplicationCommandInteractionDataOption
	data       *objects.ApplicationCommandInteractionData
}

func NewCommandOptions(options []*objects.ApplicationCommandInteractionDataOption, data *objects.ApplicationCommandInteractionData) map[string]*CommandOption {
	output := make(map[string]*CommandOption)
	for _, o := range options {
		output[o.Name] = &CommandOption{value: o.Value, options: o.Options, data: data, optionType: objects.ApplicationCommandOptionType(o.Type)}
	}
	return output
}

func (o *CommandOption) String() (string, bool) {
	if o.value == nil {
		return "", false
	}
	value, ok := o.value.(string)
	return value, ok
}

func (o *CommandOption) Integer() (int, bool) {
	if o.value == nil {
		return 0, false
	}
	value, ok := o.value.(float64)
	return int(value), ok
}

func (o *CommandOption) Boolean() (bool, bool) {
	if o.value == nil {
		return false, false
	}
	value, ok := o.value.(bool)
	return value, ok
}

func (o *CommandOption) Snowflake() (objects.Snowflake, bool) {
	if o.value == nil {
		return objects.Snowflake(0), false
	}

	intFlake, err := strconv.ParseUint(o.value.(string), 10, 64)
	if err != nil {
		return objects.Snowflake(0), false
	}

	return objects.Snowflake(intFlake), true
}

func (o *CommandOption) User() (*objects.User, bool) {
	if o.value == nil {
		return nil, false
	}

	if id, ok := o.Snowflake(); ok {
		user, ok := o.data.Resolved.Users[id]
		return &user, ok
	}

	return nil, false
}

func (o *CommandOption) Channel() (*objects.Channel, bool) {
	if o.value == nil {
		return nil, false
	}

	if id, ok := o.Snowflake(); ok {
		channel, ok := o.data.Resolved.Channels[id]
		return &channel, ok
	}

	return nil, false
}

func (o *CommandOption) Role() (*objects.Role, bool) {
	if o.value == nil {
		return nil, false
	}

	if id, ok := o.Snowflake(); ok {
		role, ok := o.data.Resolved.Roles[id]
		return &role, ok
	}

	return nil, false
}

func (o *CommandOption) SubCommand(ctx *CommandCtx) (*CommandCtx, bool) {
	if o.optionType != objects.TypeSubCommand {
		return nil, false
	}

	ctx.options = NewCommandOptions(o.options, o.data)

	return ctx, true
}
