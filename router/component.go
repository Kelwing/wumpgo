package router

import (
	"context"
	"fmt"

	"wumpgo.dev/wumpgo/objects"
)

type ComponentResponder interface {
	// UpdateMessage updates the Message the component is attached to, this is the default behavior.
	UpdateMessage() ComponentResponder
	// WithSource responses with a channel message upon function return.
	WithSource() ComponentResponder
	// Defer sends a defer response, and executes the provided function asynchronously
	Defer(ComponentHandler) ComponentResponder
	// TTS causes the content to be read out using text-to-speech
	TTS() ComponentResponder
	// Content sets the content for the response.
	Content(string) ComponentResponder
	// Contentf sets the content to the formatted value.
	Contentf(string, ...any) ComponentResponder
	// Embed adds an embed to the response, can be called up to 10 times.
	Embed(*objects.Embed) ComponentResponder
	// Embeds overwrites all embeds in the response with the provided array.
	Embeds([]*objects.Embed) ComponentResponder
	// AllowedMentions sets the allowed mentions
	AllowedMentions(*objects.AllowedMentions) ComponentResponder
	// SupressEmbeds causes all embeds on the message to be hidden
	SupressEmbeds() ComponentResponder
	// Ephemeral makes the response message ephemeral, only the person who ran the command can see it.
	Ephemeral() ComponentResponder
	// View sets the component view to respond with
	View(Renderable) ComponentResponder
	// Attach adds a file attachment to the response by Attachment ID
	Attach(f *objects.DiscordFile) ComponentResponder
	// Responds with a modal instead of a message
	// WARNING: this will override all other options you have set
	Modal(string, string, Renderable) ComponentResponder
}

var _ ComponentResponder = (*defaultComponentResponder)(nil)

func newDefaultComponentResponder() *defaultComponentResponder {
	return &defaultComponentResponder{
		response: &objects.InteractionResponse{
			Type: objects.ResponseUpdateMessage,
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

type defaultComponentResponder struct {
	response    *objects.InteractionResponse
	messageData *objects.InteractionMessagesCallbackData
	modalData   *objects.InteractionModalCallbackData
	deferFunc   ComponentHandler
	view        Renderable
	files       []*objects.DiscordFile
}

func (r *defaultComponentResponder) UpdateMessage() ComponentResponder {
	r.response.Type = objects.ResponseUpdateMessage
	return r
}

func (r *defaultComponentResponder) WithSource() ComponentResponder {
	r.response.Type = objects.ResponseChannelMessageWithSource
	return r
}

func (r *defaultComponentResponder) Defer(f ComponentHandler) ComponentResponder {
	r.deferFunc = f
	return r
}

func (r *defaultComponentResponder) TTS() ComponentResponder {
	r.messageData.TTS = true
	return r
}

func (r *defaultComponentResponder) Content(c string) ComponentResponder {
	r.messageData.Content = c
	return r
}

func (r *defaultComponentResponder) Contentf(format string, a ...any) ComponentResponder {
	r.messageData.Content = fmt.Sprintf(format, a...)
	return r
}

func (r *defaultComponentResponder) Embed(e *objects.Embed) ComponentResponder {
	r.messageData.Embeds = append(r.messageData.Embeds, e)
	return r
}

func (r *defaultComponentResponder) Embeds(e []*objects.Embed) ComponentResponder {
	r.messageData.Embeds = e
	return r
}

func (r *defaultComponentResponder) AllowedMentions(m *objects.AllowedMentions) ComponentResponder {
	r.messageData.AllowedMentions = m
	return r
}

func (r *defaultComponentResponder) SupressEmbeds() ComponentResponder {
	r.messageData.Flags |= objects.MsgFlagSupressEmbeds
	return r
}

func (r *defaultComponentResponder) Ephemeral() ComponentResponder {
	r.messageData.Flags |= objects.MsgFlagEphemeral
	return r
}

func (r *defaultComponentResponder) View(v Renderable) ComponentResponder {
	r.view = v
	return r
}

func (r *defaultComponentResponder) Attach(f *objects.DiscordFile) ComponentResponder {
	r.files = append(r.files, f)
	return r
}

func (r *defaultComponentResponder) Modal(customID string, title string, v Renderable) ComponentResponder {
	r.response.Type = objects.ResponseModal

	r.modalData = &objects.InteractionModalCallbackData{
		CustomID: customID,
		Title:    title,
	}

	return r
}

func newComponentContext(ctx context.Context, i *objects.Interaction) *ComponentContext {
	return &ComponentContext{
		InteractionContext: InteractionContext{
			interaction: i,
			ctx:         ctx,
		},
	}
}

type ComponentContext struct {
	InteractionContext
	data   *objects.MessageComponentData
	params map[string]string
}

func (c *ComponentContext) Param(name string) string {
	v, ok := c.params[name]
	if !ok {
		return ""
	}

	return v
}

func (c *ComponentContext) Values() []string {
	return c.data.Values
}

func (c *ComponentContext) SnowflakeValues() []objects.Snowflake {
	snowflakes := make([]objects.Snowflake, 0, len(c.data.Values))

	for _, v := range c.data.Values {
		s, err := objects.SnowflakeFromString(v)
		if err != nil {
			continue
		}

		snowflakes = append(snowflakes, s)
	}

	return snowflakes
}

func (c *ComponentContext) ResolveAttachment(id objects.Snowflake) *objects.Attachment {
	v, ok := c.data.Resolved.Attachments[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *ComponentContext) ResolveChannel(id objects.Snowflake) *objects.Channel {
	v, ok := c.data.Resolved.Channels[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *ComponentContext) ResolveMember(id objects.Snowflake) *objects.GuildMember {
	v, ok := c.data.Resolved.Members[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *ComponentContext) ResolveMessage(id objects.Snowflake) *objects.Message {
	v, ok := c.data.Resolved.Messages[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *ComponentContext) ResolveRole(id objects.Snowflake) *objects.Role {
	v, ok := c.data.Resolved.Roles[id]
	if !ok {
		return nil
	}

	return &v
}

func (c *ComponentContext) ResolveUser(id objects.Snowflake) *objects.User {
	v, ok := c.data.Resolved.Users[id]
	if !ok {
		return nil
	}
	return &v
}

func (c *ComponentContext) WithContext(ctx context.Context) *ComponentContext {
	c.ctx = ctx
	return c
}

type Renderable interface {
	Render() []*objects.Component
}

type View struct {
	components []Renderable
}

func NewView() *View {
	return &View{components: make([]Renderable, 0)}
}

func (v *View) Add(r ...Renderable) *View {
	v.components = append(v.components, r...)
	return v
}

func (v *View) Render() []*objects.Component {
	var components []*objects.Component

	for _, c := range v.components {
		components = append(components, c.Render()...)
	}

	return components
}

func (v *View) ToRows(maxPerRow int) []*objects.Component {
	return ComponentsToRows(v.Render(), maxPerRow)
}

func NewButton(customID ...string) *Button {
	btn := &Button{
		component: &objects.Component{
			Type:  objects.ComponentTypeButton,
			Style: objects.ButtonStylePrimary,
			Label: "Button",
		},
	}
	if len(customID) > 0 {
		btn.component.CustomID = customID[0]
	}
	return btn
}

type Button struct {
	component *objects.Component
}

func (b *Button) Style(s objects.ButtonStyle) *Button {
	b.component.Style = s
	return b
}

func (b *Button) Label(l string) *Button {
	b.component.Label = l
	return b
}

func (b *Button) URL(u string) *Button {
	b.component.Style = objects.ButtonStyleLink
	b.component.URL = u
	return b
}

func (b *Button) Emoji(e *objects.Emoji) *Button {
	b.component.Emoji = e
	return b
}

func (b *Button) Disable() *Button {
	b.component.Disabled = true
	return b
}

func (b *Button) Render() []*objects.Component {
	return []*objects.Component{b.component}
}

func NewSelectMenu(customID string) *SelectMenu {
	return &SelectMenu{
		component: &objects.Component{
			CustomID: customID,
			Type:     objects.ComponentTypeSelectMenu,
			Options:  []*objects.SelectOptions{},
		},
	}
}

type SelectMenu struct {
	component *objects.Component
}

func (s *SelectMenu) AddOption(o *objects.SelectOptions) *SelectMenu {
	s.component.Options = append(s.component.Options, o)
	return s
}

func (s *SelectMenu) Options(o []*objects.SelectOptions) *SelectMenu {
	s.component.Options = o
	return s
}

func (s *SelectMenu) ChannelTypes(t []objects.ChannelType) *SelectMenu {
	s.component.Type = objects.ComponentTypeChannelSelect
	s.component.ChannelTypes = t
	return s
}

func (s *SelectMenu) Placeholder(p string) *SelectMenu {
	s.component.Placeholder = p
	return s
}

func (s *SelectMenu) MinValues(v int) *SelectMenu {
	s.component.MinValues = &v
	return s
}

func (s *SelectMenu) MaxValues(v int) *SelectMenu {
	s.component.MaxValues = &v
	return s
}

func (s *SelectMenu) Disable() *SelectMenu {
	s.component.Disabled = true
	return s
}

func (s *SelectMenu) User() *SelectMenu {
	s.component.Type = objects.ComponentTypeUserSelect
	return s
}

func (s *SelectMenu) Role() *SelectMenu {
	s.component.Type = objects.ComponentTypeRoleSelect
	return s
}

func (s *SelectMenu) Mentionable() *SelectMenu {
	s.component.Type = objects.ComponentTypeMentionableSelect
	return s
}

func (s *SelectMenu) Render() []*objects.Component {
	return []*objects.Component{s.component}
}

func NewTextInput(customID string) *TextInput {
	return &TextInput{
		component: &objects.Component{
			CustomID: customID,
			Type:     objects.ComponentTypeInputText,
			Style:    objects.ButtonStyle(objects.TextStyleShort),
			Label:    "My Text Input",
		},
	}
}

type TextInput struct {
	component *objects.Component
}

func (t *TextInput) Short() *TextInput {
	t.component.Style = objects.ButtonStyle(objects.TextStyleShort)
	return t
}

func (t *TextInput) Long() *TextInput {
	t.component.Style = objects.ButtonStyle(objects.TextStyleParagraph)
	return t
}

func (t *TextInput) Label(l string) *TextInput {
	t.component.Label = l
	return t
}

func (t *TextInput) MinLength(v int) *TextInput {
	t.component.MinLength = &v
	return t
}

func (t *TextInput) MaxLength(v int) *TextInput {
	t.component.MaxLength = &v
	return t
}

func (t *TextInput) Required() *TextInput {
	t.component.Required = true
	return t
}

func (t *TextInput) Value(v string) *TextInput {
	t.component.Value = v
	return t
}

func (t *TextInput) Placeholder(p string) *TextInput {
	t.component.Placeholder = p
	return t
}

func (t *TextInput) Render() []*objects.Component {
	return []*objects.Component{t.component}
}
