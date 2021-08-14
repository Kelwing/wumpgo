package router

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

// ComponentRouter is used to route components.
type ComponentRouter struct {
	routes map[string]interface{}
}

// ComponentRouterCtx is used to define a components router context.
type ComponentRouterCtx struct {
	// Defines the error handler.
	errorHandler func(error) *objects.InteractionResponse

	// Defines the global allowed mentions configuration.
	globalAllowedMentions *objects.AllowedMentions

	// Defines the response builder.
	responseBuilder

	// Defines the interaction which started this.
	*objects.Interaction

	// Params are any URL params which were in the path.
	Params map[string]string `json:"params"`

	// RESTClient is used to define the REST client.
	RESTClient *rest.Client `json:"rest_client"`
}

// DeferredMessageUpdate sets the response type to DeferredMessageUpdate
// For components, ACK an interaction and edit the original message later; the user does not see a loading state
func (c *ComponentRouterCtx) DeferredMessageUpdate() *ComponentRouterCtx {
	c.respType = objects.ResponseDeferredMessageUpdate
	return c
}

// UpdateMessage sets the response type to UpdateMessage
// For components, edit the message the component was attached to
func (c *ComponentRouterCtx) UpdateMessage() *ComponentRouterCtx {
	c.respType = objects.ResponseUpdateMessage
	return c
}

// SelectMenuFunc is the function dispatched when a select menu is used.
type SelectMenuFunc func(ctx *ComponentRouterCtx, values []string) error

// Prepares the object for usage.
func (c *ComponentRouter) prep() {
	if c.routes == nil {
		c.routes = map[string]interface{}{}
	}
}

// RegisterSelectMenu is used to register a select menu route.
func (c *ComponentRouter) RegisterSelectMenu(route string, cb SelectMenuFunc) {
	c.prep()
	c.routes[route] = cb
}

// ButtonFunc is the function dispatched when a button is used.
type ButtonFunc func(ctx *ComponentRouterCtx) error

// RegisterButton is used to register a button route.
func (c *ComponentRouter) RegisterButton(route string, cb ButtonFunc) {
	c.prep()
	c.routes[route] = cb
}

// NotSelectionMenu is returned when Discord returns data that is not a selection menu.
var NotSelectionMenu = errors.New("the data returned is not that of a selection menu")

// NotButton is returned when Discord returns data that is not a button.
var NotButton = errors.New("the data returned is not that of a button")

// Adds the argument context to the handler.
type contextCallback func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string) *objects.InteractionResponse

// Used to ungeneric an error.
func ungenericError(errGeneric interface{}) error {
	var err error
	switch x := errGeneric.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = fmt.Errorf("%s", errGeneric)
	}
	return err
}

// Used to build the component router by the parent.
func (c *ComponentRouter) build(restClient *rest.Client, exceptionHandler func(err error) *objects.InteractionResponse, globalAllowedMentions *objects.AllowedMentions) interactions.HandlerFunc {
	// Build the router tree.
	c.prep()
	root := new(node)
	root.addRoute("/_postcord/void/:number", func(*objects.Interaction, *objects.ApplicationComponentInteractionData, map[string]string) *objects.InteractionResponse {
		// The point of this route is to just return nil.
		return nil
	})
	for k, v := range c.routes {
		var cb contextCallback
		switch x := v.(type) {
		case ButtonFunc:
			cb = func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string) *objects.InteractionResponse {
				if data.ComponentType != objects.ComponentTypeButton {
					return exceptionHandler(NotButton)
				}
				defer func() {
					if errGeneric := recover(); errGeneric != nil {
						// Shouldn't try and return from defer.
						exceptionHandler(ungenericError(errGeneric))
					}
				}()
				rctx := &ComponentRouterCtx{
					errorHandler:          exceptionHandler,
					globalAllowedMentions: globalAllowedMentions,
					Interaction:           ctx,
					Params:                params,
					RESTClient:            restClient,
				}
				if err := x(rctx); err != nil {
					return exceptionHandler(err)
				}
				return rctx.buildResponse(true, exceptionHandler, globalAllowedMentions)
			}
		case SelectMenuFunc:
			cb = func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string) *objects.InteractionResponse {
				values := data.Values
				if values == nil {
					// This is a blank result from Discord.
					values = []string{}
				}
				if data.ComponentType != objects.ComponentTypeSelectMenu {
					return exceptionHandler(NotSelectionMenu)
				}
				defer func() {
					if errGeneric := recover(); errGeneric != nil {
						// Shouldn't try and return from defer.
						exceptionHandler(ungenericError(errGeneric))
					}
				}()
				rctx := &ComponentRouterCtx{
					globalAllowedMentions: globalAllowedMentions,
					errorHandler:          exceptionHandler,
					Interaction:           ctx,
					Params:                params,
					RESTClient:            restClient,
				}
				if err := x(rctx, values); err != nil {
					return exceptionHandler(err)
				}
				return rctx.buildResponse(true, exceptionHandler, globalAllowedMentions)
			}
		default:
			panic("internal error: invalid interaction type")
		}
		root.addRoute(k, cb)
	}

	// Return the router.
	return func(ctx *objects.Interaction) *objects.InteractionResponse {
		params := map[string]string{}
		var data objects.ApplicationComponentInteractionData
		if err := json.Unmarshal(ctx.Data, &data); err != nil {
			return exceptionHandler(err)
		}
		route := root.getValue(data.CustomID, params)
		if route == nil {
			return nil
		}
		return route(ctx, &data, params)
	}
}
