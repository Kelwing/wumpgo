package interactions

import (
	"github.com/Postcord/objects"
)

type ButtonHandlerFunc func(ctx *ButtonCtx)

type ButtonCtx struct {
	*Ctx
	Data *objects.ApplicationComponentInteractionData
}

func (c *ButtonCtx) DeferredMessageUpdate() *ButtonCtx {
	c.Response.Type = objects.ResponseDeferredMessageUpdate
	return c
}

func (c *ButtonCtx) UpdateMessage() *ButtonCtx {
	c.Response.Type = objects.ResponseUpdateMessage
	return c
}
