package router

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type ModalHandler func(r ModalResponder, c *ModalContext)

func (r *Router) AddModalHandler(custom_id string, h ModalHandler) {
	r.modalHandlers.Insert(custom_id, h)
}

func (r *Router) routeModal(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	var data objects.ModalSubmitData
	resp := newDefaultModalResponder()

	defer func() {
		if rec := recover(); rec != nil {
			r.modalErrorHandler(resp, &errInternalCommand{rec: rec})
			response = resp.response
		}
	}()

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		return resp.response
	}

	mCtx := newModalContext(ctx, i)
	for _, c := range data.Components {
		mCtx.values[c.CustomID] = c.Value
	}

	h, ph, ok := r.modalHandlers.Search(data.CustomID)
	if !ok {
		r.modalErrorHandler(resp, ErrCustomIDNotFound)
		return resp.response
	}

	mCtx.params = ph

	h(resp, mCtx)

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

var defaultModalErrorHandler = func(r ModalResponder, err error) {
	r.Ephemeral().Content(err.Error())
}

func (r *Router) routeGatewayModal(c *rest.Client, i *objects.Interaction) {
	if i.Type != objects.InteractionModalSubmit {
		return
	}

	log.Info().Str("id", i.ID.String()).Msg("Interaction gateway event")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp := r.routeModal(ctx, i)

	log.Debug().Interface("response", resp).Msg("responding")

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create interaction response")
		r.commandErrorHandler(nil, err)
	}
}
