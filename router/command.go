package router

import (
	"context"

	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

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
	View(Renderable) CommandResponder
	// Attach adds a file attachment to the response by Attachment ID
	Attach(f *objects.DiscordFile) CommandResponder
	// Responds with a modal instead of a message
	// WARNING: this will override all other options you have set
	Modal(string, string, Renderable) CommandResponder
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
	view        Renderable
	files       []*objects.DiscordFile
}

func (r *defaultResponder) WithSource() CommandResponder {
	r.response.Type = objects.ResponseChannelMessageWithSource
	return r
}

func (r *defaultResponder) Defer(f CommandHandler) CommandResponder {
	r.response.Type = objects.ResponseDeferredChannelMessageWithSource
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

func (r *defaultResponder) View(v Renderable) CommandResponder {
	r.view = v
	return r
}

func (r *defaultResponder) Attach(f *objects.DiscordFile) CommandResponder {
	r.files = append(r.files, f)
	return r
}

func (r *defaultResponder) Modal(customID string, title string, v Renderable) CommandResponder {
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
	data        *objects.ApplicationCommandData
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

func (c *CommandContext) TargetID() objects.Snowflake {
	return c.data.TargetID
}

func (c *CommandContext) GuildID() objects.Snowflake {
	return c.interaction.GuildID
}

func (c *CommandContext) UserID() objects.Snowflake {
	if c.interaction.User != nil {
		return c.interaction.User.ID
	} else {
		return c.interaction.Member.User.ID
	}
}

func (c *CommandContext) ResolveAttachment(id objects.Snowflake) *objects.Attachment {
	v, ok := c.data.Resolved.Attachments[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *CommandContext) ResolveChannel(id objects.Snowflake) *objects.Channel {
	v, ok := c.data.Resolved.Channels[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *CommandContext) ResolveMember(id objects.Snowflake) *objects.GuildMember {
	v, ok := c.data.Resolved.Members[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *CommandContext) ResolveMessage(id objects.Snowflake) *objects.Message {
	v, ok := c.data.Resolved.Messages[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *CommandContext) ResolveRole(id objects.Snowflake) *objects.Role {
	v, ok := c.data.Resolved.Roles[id]
	if !ok {
		return nil
	}

	return &v
}

func (c *CommandContext) ResolveUser(id objects.Snowflake) *objects.User {
	v, ok := c.data.Resolved.Users[id]
	if !ok {
		return nil
	}
	return &v
}

func newCommandContext(i *objects.Interaction, opts []*objects.ApplicationCommandDataOption) *CommandContext {
	return &CommandContext{
		interaction: i,
		options:     opts,
	}
}

type CommandHandler interface {
	Handle(CommandResponder, *CommandContext)
}

type CommandHandlerFunc func(CommandResponder, *CommandContext)

func (f CommandHandlerFunc) Handle(r CommandResponder, ctx *CommandContext) {
	f(r, ctx)
}

type commandHandler struct {
	h          CommandHandler
	middleware []Middleware
}
