package interactions

import (
	"encoding/json"
	"fmt"

	"github.com/Postcord/objects"
)

type CommandRouter struct {
	handlers map[string]HandlerFunc
	app      *App
}

func NewCommandRouter(app *App) *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]HandlerFunc),
		app:      app,
	}
}

func (r *CommandRouter) AddCommand(path string, handler HandlerFunc) {
	r.handlers[path] = handler
}

func (r *CommandRouter) buildPath(prefix string, options []*objects.ApplicationCommandInteractionDataOption) (string, []*objects.ApplicationCommandInteractionDataOption) {
	finalOptions := make([]*objects.ApplicationCommandInteractionDataOption, 0)
	for _, option := range options {
		switch option.Type {
		case int(objects.TypeSubCommandGroup), int(objects.TypeSubCommand):
			return r.buildPath(fmt.Sprintf("%s/%s", prefix, option.Name), option.Options)
		default:
			finalOptions = append(finalOptions, option)
			continue
		}
	}

	return prefix, finalOptions
}

func (r *CommandRouter) Execute(c *Ctx) (err error) {
	if c.Request.Type != objects.InteractionApplicationCommand {
		err = fmt.Errorf("not a command")
		return
	}

	data := objects.ApplicationCommandInteractionData{}
	err = json.Unmarshal(c.Request.Data, &data)
	if err != nil {
		return
	}

	commandString, options := r.buildPath(data.Name, data.Options)

	if h, ok := r.handlers[commandString]; ok {
		data.Options = options
		ctx := &CommandCtx{
			Ctx:  c,
			Data: &data,
		}
		ctx.options = make(map[string]*CommandOption)
		for _, option := range data.Options {
			ctx.options[option.Name] = &CommandOption{value: option.Value, data: &data, options: option.Options, optionType: objects.ApplicationCommandOptionType(option.Type)}
		}
		h(ctx)
	} else {
		c.SetContent("Command doesn't have a handler.").Ephemeral()
		return
	}

	return
}
