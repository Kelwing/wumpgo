package router

import (
	"bytes"
	"context"
	"encoding/json"
	"runtime/debug"
	"time"

	"github.com/DataDog/gostackparse"
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/rest"
)

type ModalHandler func(r ModalResponder, c *ModalContext)

func (r *Router) AddModalHandler(custom_id string, h ModalHandler) {
	r.modalHandlers.Insert(custom_id, h)
}

func (r Router) executeModal(f func(ModalResponder, *ModalContext), cr *defaultModalResponder, ctx *ModalContext) {
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

			*cr = *newDefaultModalResponder()
			r.modalErrorHandler(cr, &errInternalCommand{rec: rec})
		}
	}()
	f(cr, ctx)
}

func (r *Router) routeModal(ctx context.Context, i *objects.Interaction) (response *objects.InteractionResponse) {
	var data objects.ModalSubmitData
	resp := newDefaultModalResponder()

	err := json.Unmarshal(i.Data, &data)
	if err != nil {
		return resp.response
	}

	mCtx := newModalContext(ctx, i)
	mCtx.client = r.client
	for _, row := range data.Components {
		for _, c := range row.Components {
			mCtx.values[c.CustomID] = ModalValue(c.Value)
		}
	}

	r.logger.Debug().Interface("values", mCtx.values).Msg("using component values")

	h, ph, ok := r.modalHandlers.Search(data.CustomID)
	if !ok {
		r.logger.Warn().Str("custom_id", data.CustomID).Msg("failed to find handler")
		r.modalErrorHandler(resp, ErrCustomIDNotFound)
		resp.response.Data = resp.messageData
		return resp.response
	}

	params := make(map[string]ModalValue)
	for k, v := range ph {
		params[k] = ModalValue(v)
	}

	mCtx.params = params

	r.executeModal(h, resp, mCtx)

	r.logger.Debug().Interface("message_data", resp.messageData).Msg("handler returned")

	if resp.view != nil {
		components := resp.view.Render()
		components = ComponentsToRows(components)

		if len(components) > 5 {
			components = components[:5]
		}

		resp.messageData.Components = components
	}

	resp.response.Data = resp.messageData

	return resp.response
}

var defaultModalErrorHandler = func(r ModalResponder, err error) {
	r.Ephemeral().Content(err.Error())
}

func (r *Router) routeGatewayModal(ctx context.Context, c *rest.Client, i *objects.InteractionCreate) {
	if i.Type != objects.InteractionModalSubmit {
		return
	}

	r.logger.Debug().Str("id", i.ID.String()).Interface("interaction", i).Msg("Interaction modal gateway event")
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	resp := r.routeModal(ctx, i.Interaction)

	r.logger.Debug().Interface("response", resp).Msg("responding")

	err := r.client.CreateInteractionResponse(ctx, i.ID, i.Token, resp)
	if err != nil {
		r.logger.Warn().Err(err).Msg("failed to create modal interaction response")
	}
}
