package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	errorHandler ErrorHandler

	// Defines the global allowed mentions configuration.
	globalAllowedMentions *objects.AllowedMentions

	// Defines the void ID generator.
	voidGenerator

	// Defines the response builder.
	responseBuilder

	// Defines the interaction which started this.
	*objects.Interaction

	// Params are any URL params which were in the path.
	Params map[string]string `json:"params"`

	// RESTClient is used to define the REST client.
	RESTClient rest.RESTClient `json:"rest_client"`
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
type contextCallback func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string, rest rest.RESTClient, errHandler ErrorHandler) *objects.InteractionResponse

// Defines the data for the context for the route.
type routeContext struct {
	cb contextCallback
	r  string
}

// Used to ungeneric an error.
func ungenericError(errGeneric interface{}) error {
	var err error
	switch x := errGeneric.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New(fmt.Sprint(errGeneric))
	}
	return err
}

// Used to build the component router by the parent.
func (c *ComponentRouter) build(loader loaderPassthrough) interactions.HandlerFunc {
	// Build the router tree.
	c.prep()
	root := new(node)
	root.addRoute("/_postcord/void/:number", &routeContext{
		cb: func(ctx *objects.Interaction, _ *objects.ApplicationComponentInteractionData, _ map[string]string, _ rest.RESTClient, _ ErrorHandler) *objects.InteractionResponse {
			// The point of this route is to just return the default handler.
			rctx := &ComponentRouterCtx{
				globalAllowedMentions: loader.globalAllowedMentions,
				Interaction:           ctx,
			}
			return rctx.buildResponse(true, nil, loader.globalAllowedMentions)
		},
		r: "/_postcord/void/:number",
	})
	for k, v := range c.routes {
		var cb contextCallback
		switch x := v.(type) {
		case ButtonFunc:
			cb = func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string, rest rest.RESTClient, errHandler ErrorHandler) *objects.InteractionResponse {
				if data.ComponentType != objects.ComponentTypeButton {
					return loader.errHandler(NotButton)
				}
				defer func() {
					if errGeneric := recover(); errGeneric != nil {
						// Shouldn't try and return from defer.
						errHandler(ungenericError(errGeneric))
					}
				}()
				rctx := &ComponentRouterCtx{
					errorHandler:          loader.errHandler,
					globalAllowedMentions: loader.globalAllowedMentions,
					Interaction:           ctx,
					Params:                params,
					RESTClient:            rest,
				}
				if err := x(rctx); err != nil {
					return errHandler(err)
				}
				return rctx.buildResponse(true, loader.errHandler, loader.globalAllowedMentions)
			}
		case SelectMenuFunc:
			cb = func(ctx *objects.Interaction, data *objects.ApplicationComponentInteractionData, params map[string]string, rest rest.RESTClient, errHandler ErrorHandler) *objects.InteractionResponse {
				values := data.Values
				if values == nil {
					// This is a blank result from Discord.
					values = []string{}
				}
				if data.ComponentType != objects.ComponentTypeSelectMenu {
					return loader.errHandler(NotSelectionMenu)
				}
				defer func() {
					if errGeneric := recover(); errGeneric != nil {
						// Shouldn't try and return from defer.
						errHandler(ungenericError(errGeneric))
					}
				}()
				rctx := &ComponentRouterCtx{
					globalAllowedMentions: loader.globalAllowedMentions,
					errorHandler:          loader.errHandler,
					Interaction:           ctx,
					Params:                params,
					RESTClient:            rest,
				}
				if err := x(rctx, values); err != nil {
					return errHandler(err)
				}
				return rctx.buildResponse(true, loader.errHandler, loader.globalAllowedMentions)
			}
		default:
			panic("postcord internal error - invalid interaction type")
		}
		root.addRoute(k, &routeContext{cb, k})
	}

	// Return the router.
	return func(ctx *objects.Interaction) *objects.InteractionResponse {
		// Create the rest tape if this is wanted.
		r := loader.rest
		tape := tape{}
		var returnedErr string
		errHandler := loader.errHandler
		if loader.generateFrames {
			r = &restTape{
				tape: &tape,
				rest: r,
			}
			errHandler = func(err error) *objects.InteractionResponse {
				returnedErr = err.Error()
				return loader.errHandler(err)
			}
		}

		// Run the command.
		params := map[string]string{}
		var data objects.ApplicationComponentInteractionData
		if err := json.Unmarshal(ctx.Data, &data); err != nil {
			return loader.errHandler(err)
		}
		route := root.getValue(data.CustomID, params)
		if route == nil {
			return nil
		}
		resp := route.cb(ctx, &data, params, r, errHandler)
		if loader.generateFrames {
			// Now we have all the data, we can generate the frame.
			fr := frame{ctx, tape, returnedErr, resp}
			go fr.write("testframes", "components", strings.ReplaceAll(route.r, "/", "_"))
		}
		return resp
	}
}
