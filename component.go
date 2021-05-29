package interactions

import (
	"github.com/Postcord/objects"
)

type ComponentHandlerFunc func(ctx *ComponentCtx)

type ComponentCtx struct {
	*Ctx
	// Data contains the component specific data payload
	Data *objects.ApplicationComponentInteractionData
}

// DeferredMessageUpdate sets the response type to DefferedMessageUpdate
// For components, ACK an interaction and edit the original message later; the user does not see a loading state
func (c *ComponentCtx) DeferredMessageUpdate() *ComponentCtx {
	c.Response.Type = objects.ResponseDeferredMessageUpdate
	return c
}

// UpdateMessage sets the response type to UpdateMessage
// For components, edit the message the component was attached to
func (c *ComponentCtx) UpdateMessage() *ComponentCtx {
	c.Response.Type = objects.ResponseUpdateMessage
	return c
}
