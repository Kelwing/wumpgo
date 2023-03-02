package dispatcher

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ Dispatcher = (*NOOPDispatcher)(nil)

type NOOPDispatcher struct {
	logger *zerolog.Logger
}

func NewNOOPDispatcher(opts ...DispatcherOption) *NOOPDispatcher {
	logger := zerolog.Nop()

	d := &NOOPDispatcher{
		logger: &logger,
	}

	for _, o := range opts {
		o(d)
	}

	return d
}

func (d *NOOPDispatcher) Dispatch(event string, data json.RawMessage) error {
	log.Trace().RawJSON("data", data).Msgf("NOOP Dispatcher: %s", event)
	return nil
}

func (d *NOOPDispatcher) setLogger(logger *zerolog.Logger) {
	d.logger = logger
}
