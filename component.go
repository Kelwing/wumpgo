package interactions

import (
	"github.com/Postcord/objects"
)

type ComponentHandlerFunc func(ctx *ComponentCtx)

type ComponentCtx struct {
	*Ctx
	Data *objects.ApplicationComponentInteractionData
}

func (c *ComponentCtx) DeferredMessageUpdate() *ComponentCtx {
	c.Response.Type = objects.ResponseDeferredMessageUpdate
	return c
}

func (c *ComponentCtx) UpdateMessage() *ComponentCtx {
	c.Response.Type = objects.ResponseUpdateMessage
	return c
}
