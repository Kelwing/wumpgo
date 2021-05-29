package interactions

import (
	"github.com/Postcord/objects"
)

type HandlerFunc func(ctx *CommandCtx)

type CommandCtx struct {
	*Ctx
	options map[string]*CommandOption
	Data    *objects.ApplicationCommandInteractionData
}

type CommandData struct {
	Command *objects.ApplicationCommand
	Handler HandlerFunc
}
