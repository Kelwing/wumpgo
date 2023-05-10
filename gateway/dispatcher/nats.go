package dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

var _ Dispatcher = (*NATSDispatcher)(nil)

type NATSDispatcher struct {
	conn   *nats.Conn
	logger *zerolog.Logger
}

func NewNATSDispatcher(url string, natsOpts []nats.Option, opts ...DispatcherOption) (*NATSDispatcher, error) {
	conn, err := nats.Connect(url, natsOpts...)
	if err != nil {
		return nil, err
	}

	logger := zerolog.Nop()

	d := &NATSDispatcher{
		conn:   conn,
		logger: &logger,
	}

	for _, o := range opts {
		o(d)
	}

	return d, nil
}

func (d *NATSDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	d.logger.Debug().Msgf("Dispatching event %s to NATS", eventName)
	return d.conn.Publish(eventName, data)
}

func (d *NATSDispatcher) SetLogger(logger *zerolog.Logger) {
	d.logger = logger
}
