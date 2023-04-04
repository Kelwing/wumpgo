package receiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/DataDog/gostackparse"
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type HandlerFunc interface{}

// Receiver is a generic interface for receiving events from a Dispatcher
type Receiver interface {
	On(handler HandlerFunc) error
	Route(event string, data json.RawMessage) error
	Run(ctx context.Context) error
}

type eventRouter struct {
	handlers   map[string][]EventHandlerIface
	log        zerolog.Logger
	client     rest.RESTClient
	errHandler func(error)
	groupName  string
}

func newEventRouter(opts ...ReceiverOption) *eventRouter {
	router := &eventRouter{
		handlers: make(map[string][]EventHandlerIface),
		log:      zerolog.Nop(),
	}

	for _, o := range opts {
		o(router)
	}

	return router
}

func (e *eventRouter) On(handler HandlerFunc) error {
	h, evt, err := eventHandlerToEvent(handler)
	if err != nil {
		return err
	}
	e.log.Debug().Str("event", evt).Msg("registered handler for event")
	e.handlers[evt] = append(e.handlers[evt], h)
	return nil
}

func (e *eventRouter) Route(event string, data json.RawMessage) error {
	defer func() {
		if rec := recover(); rec != nil {
			if e.errHandler != nil {
				e.errHandler(fmt.Errorf("%v", rec))
			}

			routines, err := gostackparse.Parse(bytes.NewReader(debug.Stack()))
			if len(err) > 0 {
				e.log.Warn().Interface("error", rec).Msg("")
			} else {
				e.log.Warn().
					Interface("error", rec).
					Str("file", routines[0].Stack[3].File).
					Int("line", routines[0].Stack[3].Line).
					Msg("")
			}

		}
	}()

	channelParts := strings.Split(strings.ToLower(event), ".")
	if len(channelParts) == 1 {
		event = channelParts[0]
	} else if len(channelParts) == 2 {
		event = channelParts[1]
	} else {
		return fmt.Errorf("invalid event name %s", event)
	}

	handlers, ok := e.handlers[event]
	if !ok {
		e.log.Debug().Msgf("received event %s, but no handlers are declared", event)
		return nil
	}

	for _, h := range handlers {
		payload := h.New()

		err := json.Unmarshal(data, payload)
		if err != nil {
			return err
		}
		ctx := context.Background()
		h.Handle(ctx, e.client, payload)
	}

	return nil
}
