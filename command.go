package interactions

import (
	"encoding/json"

	"github.com/Postcord/objects"
)

type HandlerFunc func(ctx *CommandCtx)

type CommandCtx struct {
	*Ctx
	options map[string]*CommandOption
	// Data contains the command specific data payload
	Data *objects.ApplicationCommandInteractionData
}

type CommandData struct {
	Command *objects.ApplicationCommand
	Handler HandlerFunc
}

func (c *CommandCtx) CommandName() string {
	var data objects.ApplicationCommandInteractionData
	err := json.Unmarshal(c.Request.Data, &data)
	if err != nil {
		return ""
	}
	return data.Name
}

func (c *CommandCtx) Options() []objects.ApplicationCommandInteractionDataOption {
	var data objects.ApplicationCommandInteractionData
	err := json.Unmarshal(c.Request.Data, &data)
	if err != nil {
		return nil
	}
	return data.Options
}

func (c *CommandCtx) Get(name string) *CommandOption {
	option, ok := c.options[name]
	if !ok {
		return &CommandOption{value: nil}
	}

	return option
}
