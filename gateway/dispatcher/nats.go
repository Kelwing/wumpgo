package dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type NATSDispatcher struct {
	conn *nats.Conn
}

func NewNATSDispatcher(url string, opts ...nats.Option) (*NATSDispatcher, error) {
	conn, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	return &NATSDispatcher{conn: conn}, nil
}

func (d *NATSDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	log.Debug().Msgf("Dispatching event %s to NATS", eventName)
	return d.conn.Publish(eventName, data)
}
