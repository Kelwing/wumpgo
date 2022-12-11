package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type ErrorHandler func(r CommandResponder, err error)

type CommandRouter struct {
	commandHandlers map[string]CommandHandler
	commands        []*objects.ApplicationCommand
	errHandler      ErrorHandler
	client          *rest.Client
}

func NewCommandRouter(o ...CommandRouterOption) *CommandRouter {
	r := &CommandRouter{
		commandHandlers: make(map[string]CommandHandler),
		commands:        make([]*objects.ApplicationCommand, 0),
		errHandler:      defaultErrorHandler(),
		client:          rest.New(rest.WithRateLimiter(rest.NewLeakyBucketRatelimiter())),
	}

	for _, opt := range o {
		opt(r)
	}

	return r
}

// RegisterCommand parses a struct and builds a Discord Application command from it
func (r *CommandRouter) RegisterCommand(cmd any) error {
	p := NewParser()
	cmdObj, err := p.parseCommand(reflect.ValueOf(cmd))
	if err != nil {
		return err
	}

	r.commands = append(r.commands, cmdObj)

	// Merge handlers into router handler map
	for k, v := range p.Handlers() {
		r.commandHandlers[k] = v
	}

	return nil
}

// MustRegisterCommand parses a struct and builds a Discord Application command from it.
// Panics on failure.
func (r *CommandRouter) MustRegisterCommand(cmd any) {
	if err := r.RegisterCommand(cmd); err != nil {
		panic(fmt.Sprintf("could not register command: %v", err))
	}
}

// Commands returns a slice of commands suitable for passing to
// rest.Client.BulkOverwriteGlobalCommands
func (r *CommandRouter) Commands() []*objects.ApplicationCommand {
	return r.commands
}

func (r *CommandRouter) getCommandHandler(i *objects.Interaction) (CommandHandler, []*objects.ApplicationCommandDataOption, error) {
	var data objects.ApplicationCommandData

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		return nil, nil, err
	}

	if len(data.Options) == 0 || data.Options[0].Type > objects.TypeSubCommandGroup {
		// This is definitely meant to be run as a root command
		h, ok := r.commandHandlers[data.Name]
		if ok {
			return h, data.Options, nil
		}
	}

	if data.Options[0].Type == objects.TypeSubCommand {
		h, ok := r.commandHandlers[data.Name+"/"+data.Options[0].Name]
		if !ok {
			return nil, nil, errors.New("failed to find command")
		}

		return h, data.Options[0].Options, nil
	}

	h, ok := r.commandHandlers[data.Name+"/"+data.Options[0].Name+"/"+data.Options[0].Options[0].Name]
	if !ok {
		return nil, nil, errors.New("failed to find commnad")
	}

	return h, data.Options[0].Options[0].Options, nil
}

func (r *CommandRouter) routeCommand(ctx context.Context, i *objects.Interaction) *objects.InteractionResponse {
	resp := newDefaultResponder()
	defer func() {
		if rec := recover(); rec != nil {
			r.errHandler(resp, &errInternalCommand{rec: rec})
		}
	}()
	handler, choices, err := r.getCommandHandler(i)
	if err != nil {
		r.errHandler(resp, ErrCommandNotFound)
		return resp.response
	}

	if err := unmarshalOptions(handler, choices); err != nil {
		r.errHandler(resp, ErrArgumentMismatch)
		return resp.response
	}

	cmdCtx := newCommandContext(i, choices)
	if r.client != nil {
		cmdCtx.client = r.client
	}

	handler.Handle(resp, cmdCtx)

	if resp.modalData != nil {
		resp.response.Data = resp.modalData
	} else {
		resp.response.Data = resp.messageData
	}

	return resp.response
}

func (r *CommandRouter) routeGatewayCommand(c *rest.Client, i *objects.Interaction) {
	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx := context.Background()

	if i.Type != objects.InteractionApplicationCommand {
		return
	}

	resp := r.routeCommand(ctx, i)

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
		r.errHandler(nil, err)
	}
}

func defaultErrorHandler() ErrorHandler {
	return func(r CommandResponder, err error) {
		r.Ephemeral().Content(err.Error())
	}
}
