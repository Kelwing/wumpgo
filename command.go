package interactions

import (
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
