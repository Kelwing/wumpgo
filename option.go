package interactions

import (
	"github.com/Postcord/objects"
	"github.com/mitchellh/mapstructure"
)

type CommandOption struct {
	Value interface{}
}

func (o *CommandOption) String() (string, bool) {
	if o.Value == nil {
		return "", false
	}
	value, ok := o.Value.(string)
	return value, ok
}

func (o *CommandOption) Integer() (int, bool) {
	if o.Value == nil {
		return 0, false
	}
	value, ok := o.Value.(float64)
	return int(value), ok
}

func (o *CommandOption) Boolean() (bool, bool) {
	if o.Value == nil {
		return false, false
	}
	value, ok := o.Value.(bool)
	return value, ok
}

func (o *CommandOption) User() (*objects.User, bool) {
	if o.Value == nil {
		return nil, false
	}
	user := &objects.User{}
	err := mapstructure.Decode(o.Value, user)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (o *CommandOption) Channel() (*objects.Channel, bool) {
	if o.Value == nil {
		return nil, false
	}
	channel := &objects.Channel{}
	err := mapstructure.Decode(o.Value, channel)
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (o *CommandOption) Role() (*objects.Role, bool) {
	if o.Value == nil {
		return nil, false
	}
	role := &objects.Role{}
	err := mapstructure.Decode(o.Value, role)
	if err != nil {
		return nil, false
	}
	return role, true
}
