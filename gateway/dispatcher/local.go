package dispatcher

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/rest"
)

type HandlerFunc interface{}

type LocalDispatcher struct {
	handlers map[string][]HandlerFunc
	log      zerolog.Logger
	client   *rest.Client
}

func NewLocalDispatcher(client *rest.Client, logger zerolog.Logger) *LocalDispatcher {
	return &LocalDispatcher{
		handlers: make(map[string][]HandlerFunc),
		log:      logger,
		client:   client,
	}
}

func (l *LocalDispatcher) On(evt string, handler HandlerFunc) {
	l.handlers[evt] = append(l.handlers[evt], handler)
}

func (l *LocalDispatcher) Dispatch(event string, data json.RawMessage) error {
	defer func() {
		if rec := recover(); rec != nil {
			l.log.Error().Msgf("panic while calling handler for %s", event)
		}
	}()

	handlers, ok := l.handlers[event]
	if !ok {
		l.log.Warn().Msgf("received event %s, but no handlers are declared", event)
		return nil
	}

	for _, h := range handlers {
		x := reflect.TypeOf(h)

		numIn := x.NumIn()   //Count inbound parameters
		numOut := x.NumOut() //Count outbounding parameters

		if numIn != 2 || numOut != 0 {
			l.log.Warn().Msgf("Invalid function signature for event %s. Handler: %s ", event, x.Name())
			return nil
		}

		inType := x.In(1)
		typePtr := reflect.New(inType.Elem())

		obj := typePtr.Interface()

		err := json.Unmarshal(data, obj)
		if err != nil {
			l.log.Warn().Str("event", event).Str("obj", typePtr.Type().Name()).Msgf("failed to unmarshal")
			return nil
		}

		f := reflect.ValueOf(h)
		f.Call([]reflect.Value{reflect.ValueOf(l.client), typePtr})
	}
	return nil
}
