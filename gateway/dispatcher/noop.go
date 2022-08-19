package dispatcher

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type NOOPDispatcher struct{}

func NewNOOPDispatcher() *NOOPDispatcher {
	return &NOOPDispatcher{}
}

func (d *NOOPDispatcher) Dispatch(event string, data json.RawMessage) error {
	log.Trace().RawJSON("data", data).Msgf("NOOP Dispatcher: %s", event)
	return nil
}
