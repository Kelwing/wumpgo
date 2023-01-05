package router

import (
	"context"

	"wumpgo.dev/wumpgo/objects"
)

type ModalResponder interface {
	// WithSource responses with a channel message upon function return.
	WithSource() ModalResponder
	// Defer sends a defer response, and executes the provided function asynchronously
	Defer(ModalHandler) ModalResponder
	// TTS causes the content to be read out using text-to-speech
	TTS() ModalResponder
	// Content sets the content for the response.
	Content(string) ModalResponder
	// Embed adds an embed to the response, can be called up to 10 times.
	Embed(*objects.Embed) ModalResponder
	// Embeds overwrites all embeds in the response with the provided array.
	Embeds([]*objects.Embed) ModalResponder
	// AllowedMentions sets the allowed mentions
	AllowedMentions(*objects.AllowedMentions) ModalResponder
	// SupressEmbeds causes all embeds on the message to be hidden
	SupressEmbeds() ModalResponder
	// Ephemeral makes the response message ephemeral, only the person who ran the command can see it.
	Ephemeral() ModalResponder
	// View sets the component view to respond with
	View(Renderable) ModalResponder
	// Attach adds a file attachment to the response by Attachment ID
	Attach(f *objects.DiscordFile) ModalResponder
}

var _ ModalResponder = (*defaultModalResponder)(nil)

func newDefaultModalResponder() *defaultModalResponder {
	return &defaultModalResponder{
		response: &objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
		},
		messageData: &objects.InteractionMessagesCallbackData{
			Embeds:      make([]*objects.Embed, 0),
			Components:  make([]*objects.Component, 0),
			Attachments: make([]*objects.Attachment, 0),
			Files:       make([]*objects.DiscordFile, 0),
		},
		deferFunc: nil,
		view:      nil,
		files:     make([]*objects.DiscordFile, 0),
	}
}

type defaultModalResponder struct {
	response    *objects.InteractionResponse
	messageData *objects.InteractionMessagesCallbackData
	modalData   *objects.InteractionModalCallbackData
	deferFunc   ModalHandler
	view        Renderable
	files       []*objects.DiscordFile
}

func (r *defaultModalResponder) WithSource() ModalResponder {
	r.response.Type = objects.ResponseChannelMessageWithSource
	return r
}

func (r *defaultModalResponder) Defer(f ModalHandler) ModalResponder {
	r.response.Type = objects.ResponseDeferredChannelMessageWithSource
	r.deferFunc = f
	return r
}

func (r *defaultModalResponder) TTS() ModalResponder {
	r.messageData.TTS = true
	return r
}

func (r *defaultModalResponder) Content(c string) ModalResponder {
	r.messageData.Content = c
	return r
}

func (r *defaultModalResponder) Embed(e *objects.Embed) ModalResponder {
	r.messageData.Embeds = append(r.messageData.Embeds, e)
	return r
}

func (r *defaultModalResponder) Embeds(e []*objects.Embed) ModalResponder {
	r.messageData.Embeds = e
	return r
}

func (r *defaultModalResponder) AllowedMentions(m *objects.AllowedMentions) ModalResponder {
	r.messageData.AllowedMentions = m
	return r
}

func (r *defaultModalResponder) SupressEmbeds() ModalResponder {
	r.messageData.Flags |= objects.MsgFlagSupressEmbeds
	return r
}

func (r *defaultModalResponder) Ephemeral() ModalResponder {
	r.messageData.Flags |= objects.MsgFlagEphemeral
	return r
}

func (r *defaultModalResponder) View(v Renderable) ModalResponder {
	r.view = v
	return r
}

func (r *defaultModalResponder) Attach(f *objects.DiscordFile) ModalResponder {
	r.files = append(r.files, f)
	return r
}

func newModalContext(ctx context.Context, i *objects.Interaction) *ModalContext {
	return &ModalContext{
		values: make(map[string]string),
	}
}

type ModalContext struct {
	InteractionContext
	values map[string]string
	params map[string]string
}

// Value returns the value of a text input component within the modal
func (c *ModalContext) Value(customID string) string {
	v, ok := c.values[customID]
	if !ok {
		return ""
	}

	return v
}
