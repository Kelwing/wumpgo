package router

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/Postcord/objects"
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

// UnsetModalRouter is thrown when the modal router is unset.
var UnsetModalRouter = errors.New("modal router is unset")

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

// Internal struct to expose the response builder whilst having a memory friendly way to change the generic.
// !! UNSAFE WARNING !!: This uses the location of the struct to determine the where the parent is. This is fine
// because it is embedded into the struct as the first item. HOWEVER, IF YOU IMPLEMENT THIS, YOU MUST MAKE IT THE
// FIRST ITEM IN THE STRUCT OR YOU WILL WRITE TO RANDOM MEMORY. BE CAREFUL! The reason this is unsafe in the first
// place is to avoid unneeded circular references.
type publicResponseBuilder[T any] struct {
	responseBuilder
}

// See warning above.
func (c *publicResponseBuilder[T]) getOrigin() T {
	unsafePtr := unsafe.Pointer(c)
	var ptr *T
	switch (any)(ptr).(type) {
	case **ModalRouterCtx:
		x := (*ModalRouterCtx)(unsafePtr)
		return (any)(x).(T)
	case **CommandRouterCtx:
		x := (*CommandRouterCtx)(unsafePtr)
		return (any)(x).(T)
	case **ComponentRouterCtx:
		x := (*ComponentRouterCtx)(unsafePtr)
		return (any)(x).(T)
	default:
		panic("postcord internal error - unknown type of parent for public response builder")
	}
}

// SetEmbed is used to set the embed, overwriting any previously.
func (c *publicResponseBuilder[T]) SetEmbed(embed *objects.Embed) T {
	c.editEmbed(embed, false)
	return c.getOrigin()
}

// AddEmbed is used to append the embed, joining any previously.
func (c *publicResponseBuilder[T]) AddEmbed(embed *objects.Embed) T {
	c.editEmbed(embed, true)
	return c.getOrigin()
}

// AddComponentRow is used to add a row of components.
func (c *publicResponseBuilder[T]) AddComponentRow(row []*objects.Component) T {
	component := &objects.Component{Type: objects.ComponentTypeActionRow, Components: row}
	response := c.ResponseData()
	response.Components = append(response.Components, component)
	return c.getOrigin()
}

// SetComponentRows is used to set rows of components.
func (c *publicResponseBuilder[T]) SetComponentRows(rows [][]*objects.Component) T {
	components := make([]*objects.Component, len(rows))
	for i, v := range rows {
		components[i] = &objects.Component{Type: objects.ComponentTypeActionRow, Components: v}
	}
	c.ResponseData().Components = components
	return c.getOrigin()
}

// ClearComponents is used to clear the components in a response.
func (c *publicResponseBuilder[T]) ClearComponents() T {
	c.ResponseData().Components = []*objects.Component{}
	return c.getOrigin()
}

// SetContent is used to set the content of a response.
func (c *publicResponseBuilder[T]) SetContent(content string) T {
	c.ResponseData().Content = content
	return c.getOrigin()
}

// SetContentf is used to set the content of a response using fmt.Sprintf.
func (c *publicResponseBuilder[T]) SetContentf(content string, args ...any) T {
	c.ResponseData().Content = fmt.Sprintf(content, args...)
	return c.getOrigin()
}

// SetAllowedMentions is used to set the allowed mentions of a response. This will override your global configuration.
func (c *publicResponseBuilder[T]) SetAllowedMentions(config *objects.AllowedMentions) T {
	c.ResponseData().AllowedMentions = config
	return c.getOrigin()
}

// SetTTS is used to set the TTS configuration for your response.
func (c *publicResponseBuilder[T]) SetTTS(tts bool) T {
	c.ResponseData().TTS = tts
	return c.getOrigin()
}

// Ephemeral is used to set the response as ephemeral.
func (c *publicResponseBuilder[T]) Ephemeral() T {
	c.ResponseData().Flags = 64
	return c.getOrigin()
}

// AttachBytes adds a file attachment to the response from a byte array
func (c *publicResponseBuilder[T]) AttachBytes(data []byte, filename, description string) T {
	file := &objects.DiscordFile{
		Buffer:      bytes.NewBuffer(data),
		Filename:    filename,
		Description: description,
	}
	response := c.ResponseData()
	response.Files = append(response.Files, file)
	return c.getOrigin()
}

// AttachFile adds a file attachment to the response from an *objects.DiscordFile
func (c *publicResponseBuilder[T]) AttachFile(file *objects.DiscordFile) T {
	response := c.ResponseData()
	response.Files = append(response.Files, file)
	return c.getOrigin()
}

// ChannelMessageWithSource is used to respond to the interaction with a message.
func (c *publicResponseBuilder[T]) ChannelMessageWithSource() T {
	c.respType = objects.ResponseChannelMessageWithSource
	return c.getOrigin()
}
