package receiver

import (
	"context"

	"github.com/nats-io/nats.go"
)

var _ Receiver = (*NATSReceiver)(nil)

type NATSReceiver struct {
	*eventRouter
	conn *nats.Conn
}

func NewNATSReceiver(url string, natsOptions []nats.Option, opts ...ReceiverOption) (*NATSReceiver, error) {
	conn, err := nats.Connect(url, natsOptions...)
	if err != nil {
		return nil, err
	}

	router := newEventRouter(opts...)

	return &NATSReceiver{conn: conn, eventRouter: router}, nil
}

func (r *NATSReceiver) Run(ctx context.Context) error {
	ch := make(chan *nats.Msg, 64)
	var sub *nats.Subscription
	var err error
	if r.groupName != "" {
		sub, err = r.conn.QueueSubscribeSyncWithChan("discord.*", r.groupName, ch)
	} else {
		sub, err = r.conn.ChanSubscribe("discord.*", ch)
	}
	if err != nil {
		return nil
	}
	defer sub.Unsubscribe()
	for {
		select {
		case msg := <-ch:
			if err := r.Route(msg.Subject, msg.Data); err != nil {
				r.log.Warn().Err(err).Str("event", msg.Subject).Msg("failed to route event")
			}
		case <-ctx.Done():
			return nil
		}
	}
}
