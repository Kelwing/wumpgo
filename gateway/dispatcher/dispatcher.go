package dispatcher

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

type Dispatcher interface {
	Dispatch(event string, data json.RawMessage) error
	setLogger(l *zerolog.Logger)
}

type DispatcherOption func(d Dispatcher)

func WithLogger(l *zerolog.Logger) DispatcherOption {
	return func(d Dispatcher) {
		d.setLogger(l)
	}
}
