package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/DataDog/gostackparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type Middleware func(next CommandHandlerFunc) CommandHandlerFunc

// RegisterCommand parses a struct and builds a Discord Application command from it
func (r *Router) RegisterCommand(cmd any, m ...Middleware) error {
	p := NewParser()
	cmdObj, err := p.parseCommand(reflect.ValueOf(cmd))
	if err != nil {
		r.logger.Warn().Err(err).Msg("failed to add command")
		return err
	}

	r.commands = append(r.commands, cmdObj)

	if len(p.Handlers()) == 0 {
		r.logger.Warn().Msgf("no handlers registered for command %s", cmdObj.Name)
	}

	// Merge handlers into router handler map
	for k, v := range p.Handlers() {
		r.commandHandlers[k] = &commandHandler{
			h:          v,
			middleware: m,
		}
	}

	return nil
}

// MustRegisterCommand parses a struct and builds a Discord Application command from it.
// Panics on failure.
func (r *Router) MustRegisterCommand(cmd any, m ...Middleware) {
	if err := r.RegisterCommand(cmd, m...); err != nil {
		panic(fmt.Errorf("could not register command: %v", err))
	}
}

// Commands returns a slice of commands suitable for passing to
// rest.Client.BulkOverwriteGlobalCommands
func (r *Router) Commands() []*objects.ApplicationCommand {
	return r.commands
}

func (r *Router) getCommandHandler(data *objects.ApplicationCommandData) (*commandHandler, []*objects.ApplicationCommandDataOption, error) {
	r.logger.Debug().Interface("data", data).Interface("handlers", r.commandHandlers).Msg("getting command handler")
	if len(data.Options) == 0 || data.Options[0].Type > objects.TypeSubCommandGroup {
		// This is definitely meant to be run as a root command
		h, ok := r.commandHandlers[data.Name]
		if ok {
			return h, data.Options, nil
		} else {
			return nil, nil, fmt.Errorf("command not registered")
		}
	}

	if data.Options[0].Type == objects.TypeSubCommand {
		h, ok := r.commandHandlers[data.Name+"/"+data.Options[0].Name]
		if !ok {
			return nil, nil, errors.New("command not registered")
		}

		return h, data.Options[0].Options, nil
	}

	h, ok := r.commandHandlers[data.Name+"/"+data.Options[0].Name+"/"+data.Options[0].Options[0].Name]
	if !ok {
		return nil, nil, errors.New("failed to find commnad")
	}

	return h, data.Options[0].Options[0].Options, nil
}

func (r Router) executeCommand(f func(CommandResponder, *CommandContext), cr *defaultResponder, ctx *CommandContext) {
	defer func() {
		if rec := recover(); rec != nil {
			routines, errs := gostackparse.Parse(bytes.NewReader(debug.Stack()))
			if len(errs) > 0 {
				r.logger.Warn().Interface("error", rec).Msg("")
			} else {
				arr := zerolog.Arr()
				for _, f := range routines[0].Stack {
					arr.Interface(f)
				}
				r.logger.Warn().
					Interface("error", rec).
					Array("stack", arr).
					Msg("")
			}

			*cr = *newDefaultResponder()
			r.commandErrorHandler(cr, &errInternalCommand{rec: rec})
		}
	}()
	f(cr, ctx)
}

func (r *Router) routeCommand(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	resp := newDefaultResponder()

	var data objects.ApplicationCommandData

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		r.commandErrorHandler(resp, fmt.Errorf("not a command interaction"))
		return resp.response
	}

	handler, choices, err := r.getCommandHandler(&data)
	if err != nil {
		r.commandErrorHandler(resp, ErrCommandNotFound)
		return resp.response
	}

	handlerType := reflect.TypeOf(handler.h)
	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
	}
	newHandler := reflect.New(handlerType).Interface().(CommandHandler)

	if err := unmarshalOptions(newHandler, choices, &data.Resolved, r.logger); err != nil {
		r.commandErrorHandler(resp, ErrArgumentMismatch)
		return resp.response
	}

	cmdCtx := newCommandContext(ctx, i, choices)
	if r.client != nil {
		cmdCtx.client = r.client
	}

	cmdCtx.data = &data

	h := newHandler.Handle

	log.Debug().Msgf("chaining %d middleware", len(handler.middleware))

	for i := len(handler.middleware) - 1; i >= 0; i-- {
		h = handler.middleware[i](h)
	}

	r.executeCommand(h, resp, cmdCtx)

	if resp.view != nil {
		components := resp.view.Render()

		if resp.response.Type != objects.ResponseModal {
			components = ComponentsToRows(components)
		} else {
			components = ComponentsToRows(components, 1)
		}

		if len(components) > 5 {
			components = components[:5]
		}

		if resp.response.Type == objects.ResponseModal {
			resp.modalData.Components = components
		} else {
			resp.messageData.Components = components
		}
	}

	if resp.modalData != nil {
		resp.response.Data = resp.modalData
	} else {
		resp.response.Data = resp.messageData
	}

	if resp.deferFunc != nil {
		go func() {
			log.Debug().Msg("Starting defer handler")
			dr := newDefaultResponder()
			resp.deferFunc.Handle(dr, cmdCtx)
			log.Debug().Msg("Defer handler finished")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_, err := r.client.EditOriginalInteractionResponse(ctx, i.ApplicationID, i.Token, &rest.EditWebhookMessageParams{
				Content:         dr.messageData.Content,
				Embeds:          dr.messageData.Embeds,
				AllowedMentions: dr.messageData.AllowedMentions,
				Components:      dr.messageData.Components,
			})
			if err != nil {
				log.Warn().Err(err).Msg("failed to update interaction")
			}
		}()
	}

	return resp.response
}

func (r *Router) routeGatewayCommand(ctx context.Context, c *rest.Client, i *objects.InteractionCreate) {
	if i.Type != objects.InteractionApplicationCommand {
		return
	}

	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp := r.routeCommand(ctx, i.Interaction)

	log.Debug().Interface("response", resp).Msg("responding")

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
	}
}

var defaultCommandErrorHandler = func(r CommandResponder, err error) {
	r.Ephemeral().Content(err.Error())
}

func (r *Router) routeAutocomplete(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	var data objects.ApplicationCommandData

	defer func() {
		if rec := recover(); rec != nil {
			response = &objects.InteractionResponse{
				Type: objects.ResponseCommandAutocompleteResult,
				Data: &objects.InteractionAutocompleteCallbackData{
					Choices: []*objects.ApplicationCommandOptionChoice{},
				},
			}
		}
	}()

	respData := &objects.InteractionAutocompleteCallbackData{
		Choices: []*objects.ApplicationCommandOptionChoice{},
	}

	resp := &objects.InteractionResponse{
		Type: objects.ResponseCommandAutocompleteResult,
		Data: respData,
	}

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		return resp
	}

	handler, choices, err := r.getCommandHandler(&data)
	if err != nil {
		return resp
	}

	ac, ok := handler.h.(AutoCompleter)
	if !ok {
		return resp
	}

	if err := unmarshalOptions(handler.h, choices, &data.Resolved); err != nil {
		return resp
	}

	var name string
	var value interface{}

	for _, o := range data.Options {
		if o.Focused {
			name = o.Name
			value = o.Value
			break
		}
	}

	respData.Choices = ac.AutoComplete(name, value)
	return resp
}

func (r *Router) routeGatewayAutocomplete(ctx context.Context, c *rest.Client, i *objects.InteractionCreate) {
	if i.Type != objects.InteractionAutoComplete {
		return
	}

	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp := r.routeAutocomplete(ctx, i.Interaction)

	log.Debug().Interface("response", resp).Msg("responding")

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
	}
}
