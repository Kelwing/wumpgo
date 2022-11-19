package dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

type NSQDispatcher struct {
	conn *nsq.Producer
}

func NewNSQDispatcher(url string, config *nsq.Config) (*NSQDispatcher, error) {
	conn, err := nsq.NewProducer(url, config)
	if err != nil {
		return nil, err
	}

	return &NSQDispatcher{conn: conn}, nil
}

func (d *NSQDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	log.Debug().Msgf("Dispatching event %s to NSQ", eventName)
	return d.conn.Publish(eventName, data)
}
