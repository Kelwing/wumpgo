package router

import (
	"context"

	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type View interface{}

type CommandResponder interface {
	// WithSource responses with a channel message upon function return, this is the default behavior.
	WithSource() CommandResponder
	// Defer sends a defer response, and executes the provided function asynchronously
	Defer(CommandHandler) CommandResponder
	// TTS causes the content to be read out using text-to-speech
	TTS() CommandResponder
	// Content sets the content for the response.
	Content(string) CommandResponder
	// Embed adds an embed to the response, can be called up to 10 times.
	Embed(*objects.Embed) CommandResponder
	// Embeds overwrites all embeds in the response with the provided array.
	Embeds([]*objects.Embed) CommandResponder
	// AllowedMentions sets the allowed mentions
	AllowedMentions(*objects.AllowedMentions) CommandResponder
	// SupressEmbeds causes all embeds on this message to be hidden
	SupressEmbeds() CommandResponder
	// Ephemeral makes the response message ephemeral, only the person who ran the command can see it.
	Ephemeral() CommandResponder
	// View sets the component view to respond with
	View(View) CommandResponder
	// Attach adds a file attachment to the response by Attachment ID
	Attach(f *objects.DiscordFile) CommandResponder
	// Responds with a modal instead of a message
	// WARNING: this will override all other options you have set
	Modal(string, string, View) CommandResponder
}

func newDefaultResponder() *defaultResponder {
	return &defaultResponder{
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

type defaultResponder struct {
	response    *objects.InteractionResponse
	messageData *objects.InteractionMessagesCallbackData
	modalData   *objects.InteractionModalCallbackData
	deferFunc   CommandHandler
	view        View
	files       []*objects.DiscordFile
}

func (r *defaultResponder) WithSource() CommandResponder {
	r.response.Type = objects.ResponseChannelMessageWithSource
	return r
}

func (r *defaultResponder) Defer(f CommandHandler) CommandResponder {
	r.deferFunc = f
	return r
}

func (r *defaultResponder) TTS() CommandResponder {
	r.messageData.TTS = true
	return r
}

func (r *defaultResponder) Content(c string) CommandResponder {
	r.messageData.Content = c
	return r
}

func (r *defaultResponder) Embed(e *objects.Embed) CommandResponder {
	r.messageData.Embeds = append(r.messageData.Embeds, e)
	return r
}

func (r *defaultResponder) Embeds(e []*objects.Embed) CommandResponder {
	r.messageData.Embeds = e
	return r
}

func (r *defaultResponder) AllowedMentions(m *objects.AllowedMentions) CommandResponder {
	r.messageData.AllowedMentions = m
	return r
}

func (r *defaultResponder) SupressEmbeds() CommandResponder {
	r.messageData.Flags |= objects.MsgFlagSupressEmbeds
	return r
}

func (r *defaultResponder) Ephemeral() CommandResponder {
	r.messageData.Flags |= objects.MsgFlagEphemeral
	return r
}

func (r *defaultResponder) View(v View) CommandResponder {
	r.view = v
	return r
}

func (r *defaultResponder) Attach(f *objects.DiscordFile) CommandResponder {
	r.files = append(r.files, f)
	return r
}

func (r *defaultResponder) Modal(customID string, title string, v View) CommandResponder {
	r.response.Type = objects.ResponseModal

	r.modalData = &objects.InteractionModalCallbackData{
		CustomID: customID,
		Title:    title,
	}

	return r
}

type CommandContext struct {
	interaction *objects.Interaction
	options     []*objects.ApplicationCommandDataOption
	ctx         context.Context
	client      *rest.Client
}

func (c *CommandContext) Interaction() *objects.Interaction {
	return c.interaction
}

func (c *CommandContext) Options() []*objects.ApplicationCommandDataOption {
	return c.options
}

func (c *CommandContext) Context() context.Context {
	return c.ctx
}

func (c *CommandContext) WithContext(ctx context.Context) *CommandContext {
	c.ctx = ctx
	return c
}

func (c *CommandContext) Client() *rest.Client {
	return c.client
}

func newCommandContext(i *objects.Interaction, opts []*objects.ApplicationCommandDataOption) *CommandContext {
	return &CommandContext{
		interaction: i,
		options:     opts,
	}
}

// CommandHandler defines a command handler contract
type CommandHandler interface {
	Handle(CommandResponder, *CommandContext)
}

// CommandHandlerFunc is an adapter to allow using functions
// as command handlers.  Useful for CommandRouter.SetPreHandler
// and CommandResponder.Defer
type CommandHandlerFunc func(CommandResponder, *CommandContext)

// Handle implements the CommandHandler interface
func (f CommandHandlerFunc) Handle(r CommandResponder, ctx *CommandContext) {
	f(r, ctx)
}

// CommandMiddleware defines a middleware structure to wrap command handlers
type CommandMiddleware func(next CommandHandler) CommandHandler

func defaultMiddleware(next CommandHandler) CommandHandler {
	return CommandHandlerFunc(func(r CommandResponder, ctx *CommandContext) {
		next.Handle(r, ctx)
	})
}
