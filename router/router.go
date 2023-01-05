package router

import (
	"context"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type CommandErrorHandler func(r CommandResponder, err error)
type ComponentErrorHandler func(r ComponentResponder, err error)
type ModalErrorHandler func(r ModalResponder, err error)

type Router struct {
	commandHandlers       map[string]CommandHandler
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
		commandHandlers:       make(map[string]CommandHandler),
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

func (c *InteractionContext) Interaction() *objects.Interaction {
	return c.interaction
}

func (c *InteractionContext) Client() *rest.Client {
	return c.client
}

func (c *InteractionContext) Message() *objects.Message {
	return c.interaction.Message
}

func (c *InteractionContext) Member() *objects.GuildMember {
	return c.interaction.Member
}

func (c *InteractionContext) User() *objects.User {
	return c.interaction.User
}

func (c *InteractionContext) Context() context.Context {
	return c.ctx
}

func (c *InteractionContext) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}
