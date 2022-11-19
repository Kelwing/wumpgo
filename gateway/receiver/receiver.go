package receiver

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type HandlerFunc interface{}

// Receiver is a generic interface for receiving events from a Dispatcher
type Receiver interface {
	On(evt string, handler HandlerFunc)
}

type eventRouter struct {
	handlers   map[string][]HandlerFunc
	log        zerolog.Logger
	client     *rest.Client
	errHandler func(error)
}

func newEventRouter(opts ...ReceiverOption) *eventRouter {
	router := &eventRouter{
		handlers: make(map[string][]HandlerFunc),
		log:      zerolog.Nop(),
	}

	for _, o := range opts {
		o(router)
	}

	return router
}

func (e *eventRouter) On(evt string, handler HandlerFunc) {
	evt = strings.ToLower(evt)
	e.handlers[evt] = append(e.handlers[evt], handler)
}

func (e *eventRouter) Route(event string, data json.RawMessage) error {
	defer func() {
		if rec := recover(); rec != nil {
			if e.errHandler != nil {
				e.errHandler(fmt.Errorf("%v", rec))
			}

			e.log.Warn().Stack().Interface("error", rec).Msg("")
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
		e.log.Warn().Msgf("received event %s, but no handlers are declared", event)
		return nil
	}

	for _, h := range handlers {
		x := reflect.TypeOf(h)

		numIn := x.NumIn()
		numOut := x.NumOut()

		if numIn != 2 || numOut != 0 {
			e.log.Warn().Msgf("Invalid function signature for event %s. Handler: %s ", event, x.Name())
			return nil
		}

		inType := x.In(1)
		typePtr := reflect.New(inType.Elem())

		obj := typePtr.Interface()

		err := json.Unmarshal(data, obj)
		if err != nil {
			e.log.Warn().Str("event", event).Str("obj", typePtr.Type().Name()).Msgf("failed to unmarshal")
			return nil
		}

		f := reflect.ValueOf(h)
		f.Call([]reflect.Value{reflect.ValueOf(e.client), typePtr})
	}
	return nil
}
