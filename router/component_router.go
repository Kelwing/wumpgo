package router

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type ComponentHandler func(r ComponentResponder, c *ComponentContext)

func (r *Router) AddHandler(custom_id string, h ComponentHandler) {
	r.componentHandlers.Insert(custom_id, h)
}

func (r *Router) routeComponent(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	var data objects.MessageComponentData
	resp := newDefaultComponentResponder()

	defer func() {
		if rec := recover(); rec != nil {
			r.componentErrorHandler(resp, &errInternalCommand{rec: rec})
			response = resp.response
		}
	}()

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

	h(resp, cmpCtx)

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

func (r *Router) routeGatewayComponent(c *rest.Client, i *objects.Interaction) {
	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx := context.Background()

	if i.Type != objects.InteractionComponent {
		return
	}

	resp := r.routeComponent(ctx, i)

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
		r.componentErrorHandler(nil, err)
	}
}

var defaultComponentErrorHandler = func(r ComponentResponder, err error) {
	r.Ephemeral().Content(err.Error())
}
