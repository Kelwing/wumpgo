package router

import (
	"errors"
	"github.com/Postcord/objects"
	"sync"
)

// Embedded into contexts to build responses.
type responseBuilder struct {
	// Defines the response type. If this is zero, it is inferred from the content.
	respType objects.ResponseType

	// Defines the data pointer.
	dataPtr     *objects.InteractionApplicationCommandCallbackData
	dataPtrLock sync.Mutex
}

// ResponseData is used to return a pointer to the response data. The data will be created if it doesn't exist, so it'll never be nil.
// NOTE: If the type is being inferred, this will mark this as a message update. Use the helper function for the type you want to prevent this.
func (r *responseBuilder) ResponseData() *objects.InteractionApplicationCommandCallbackData {
	r.dataPtrLock.Lock()
	x := r.dataPtr
	if x == nil {
		x = &objects.InteractionApplicationCommandCallbackData{}
		r.dataPtr = x
	}
	r.dataPtrLock.Unlock()
	return x
}

// NoCommandResponse is thrown when the application doesn't respond for a command.
var NoCommandResponse = errors.New("expected data for command response")

// Builds the response.
func (r *responseBuilder) buildResponse(component bool, errorHandler ErrorHandler, globalAllowedMentions *objects.AllowedMentions) *objects.InteractionResponse {
	// Get the content and do not try and create it.
	r.dataPtrLock.Lock()
	data := r.dataPtr
	r.dataPtrLock.Unlock()

	// Get the response type.
	respType := r.respType
	if respType == 0 {
		// We should try and infer the type.
		if data == nil {
			if !component {
				// If this isn't a component, something has gone badly wrong.
				return errorHandler(NoCommandResponse)
			}

			// The response type is deferred message update.
			respType = objects.ResponseDeferredMessageUpdate
		} else if component {
			// The type is message update.
			respType = objects.ResponseUpdateMessage
		} else {
			// The type is message create.
			respType = objects.ResponseChannelMessageWithSource
		}
	}

	// Handle global allowed mentions.
	if data != nil && data.AllowedMentions == nil {
		data.AllowedMentions = globalAllowedMentions
	}

	// Create the object.
	return &objects.InteractionResponse{
		Type: respType,
		Data: data,
	}
}

// Internal method to edit embeds.
func (r *responseBuilder) editEmbed(embed *objects.Embed, appendEmbed bool) {
	if embed == nil {
		return
	}
	d := r.ResponseData()
	if appendEmbed {
		d.Embeds = append(d.Embeds, embed)
	} else {
		d.Embeds = []*objects.Embed{embed}
	}
}

// Internal method to edit components.
func (r *responseBuilder) editComponent(component *objects.Component, appendComponent bool) {
	if component == nil {
		return
	}
	d := r.ResponseData()
	if appendComponent {
		d.Components = append(d.Components, component)
	} else {
		d.Components = []*objects.Component{component}
	}
}
