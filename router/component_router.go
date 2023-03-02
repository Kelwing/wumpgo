package router

import (
	"bytes"
	"context"
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/DataDog/gostackparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type ComponentHandler func(r ComponentResponder, c *ComponentContext)

func (r *Router) AddHandler(custom_id string, h ComponentHandler) {
	r.componentHandlers.Insert(custom_id, h)
}

func (r Router) executeComponent(f func(ComponentResponder, *ComponentContext), cr *defaultComponentResponder, ctx *ComponentContext) {
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

			*cr = *newDefaultComponentResponder()
			r.componentErrorHandler(cr, &errInternalCommand{rec: rec})
		}
	}()
	f(cr, ctx)
}

func (r *Router) routeComponent(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	var data objects.MessageComponentData
	resp := newDefaultComponentResponder()

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		return resp.response
	}

	cmpCtx := newComponentContext(ctx, i)

	cmpCtx.data = &data

	h, ph, ok := r.componentHandlers.Search(data.CustomID)
	if !ok {
		r.componentErrorHandler(resp, ErrCustomIDNotFound)
		return resp.response
	}

	cmpCtx.params = ph

	r.executeComponent(h, resp, cmpCtx)

	if resp.view != nil {
		components := resp.view.Render()

		if resp.response.Type != objects.ResponseModal {
			components = ComponentsToRows(components)
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

	return resp.response
}

func (r *Router) routeGatewayComponent(ctx context.Context, c *rest.Client, i *objects.InteractionCreate) {
	if i.Type != objects.InteractionComponent {
		return
	}

	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp := r.routeComponent(ctx, i.Interaction)

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
		r.componentErrorHandler(nil, err)
	}
}

var defaultComponentErrorHandler = func(r ComponentResponder, err error) {
	r.Ephemeral().Content(err.Error())
}
