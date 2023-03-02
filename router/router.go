package router

import (
	"context"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/objects/permissions"
	"wumpgo.dev/wumpgo/rest"
)

type CommandErrorHandler func(r CommandResponder, err error)
type ComponentErrorHandler func(r ComponentResponder, err error)
type ModalErrorHandler func(r ModalResponder, err error)

type Router struct {
	commandHandlers       map[string]*commandHandler
	commands              []*objects.ApplicationCommand
	componentErrorHandler ComponentErrorHandler
	commandErrorHandler   CommandErrorHandler
	modalErrorHandler     ModalErrorHandler
	client                *rest.Client
	componentHandlers     *trieNode[ComponentHandler]
	modalHandlers         *trieNode[ModalHandler]
	logger                zerolog.Logger
}

func New(opts ...RouterOption) *Router {
	r := &Router{
		commandHandlers:       make(map[string]*commandHandler),
		commands:              make([]*objects.ApplicationCommand, 0),
		componentErrorHandler: defaultComponentErrorHandler,
		commandErrorHandler:   defaultCommandErrorHandler,
		modalErrorHandler:     defaultModalErrorHandler,
		client:                rest.New(rest.WithRateLimiter(rest.NewLeakyBucketRatelimiter())),
		componentHandlers:     newTrieNode[ComponentHandler](),
		modalHandlers:         newTrieNode[ModalHandler](),
		logger:                zerolog.Nop(),
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

type Context interface {
	Interaction() *objects.Interaction
	Context() context.Context
	WithContext(ctx context.Context) Context
	Client() *rest.Client
	Message() *objects.Message
	Member() *objects.GuildMember
	User() *objects.User
}

var _ Context = (*InteractionContext)(nil)

type InteractionContext struct {
	interaction *objects.Interaction
	ctx         context.Context
	client      *rest.Client
}

// Interaction payload
func (c *InteractionContext) Interaction() *objects.Interaction {
	return c.interaction
}

// Client is a Discord REST client, nil if not configured
func (c *InteractionContext) Client() *rest.Client {
	return c.client
}

// Message object for the interaction
func (c *InteractionContext) Message() *objects.Message {
	return c.interaction.Message
}

// Member who invoked the interaction, nil if not invoked in a Guild
func (c *InteractionContext) Member() *objects.GuildMember {
	return c.interaction.Member
}

// User who invoked the interaction, guaranteed to not be nil
func (c *InteractionContext) User() *objects.User {
	if c.interaction.Member == nil {
		return c.interaction.User
	} else {
		return c.interaction.Member.User
	}
}

// Context for the interaction, canceled when the interaction is no longer valid
func (c *InteractionContext) Context() context.Context {
	return c.ctx
}

// ChannelID of the channel the interaction was invoked in
func (c *InteractionContext) ChannelID() objects.Snowflake {
	return c.interaction.ChannelID
}

// GuildID of the guild the interaction was invoked in, 0 if not invoked in a Guild
func (c *InteractionContext) GuildID() objects.Snowflake {
	return c.interaction.GuildID
}

// Locale of the user who invoked the interaction
func (c *InteractionContext) Locale() string {
	return c.interaction.Locale
}

// Prefered locale for the guild the interaction was invoked in, empty string if not invoked in a Guild
func (c *InteractionContext) GuildLocale() string {
	return c.interaction.GuildLocale
}

// Permissions of the app in the context the interaction was invoked in
func (c *InteractionContext) Permissions() permissions.PermissionBit {
	return c.interaction.AppPermissions
}

// AuthorID is the UserID of the user that invoked the interaction
func (c *InteractionContext) AuthorID() objects.Snowflake {
	if c.interaction.User != nil {
		return c.interaction.User.ID
	} else {
		return c.interaction.Member.User.ID
	}
}

// WithContext attaches a new context.Context to this InteractionContext
func (c *InteractionContext) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}
