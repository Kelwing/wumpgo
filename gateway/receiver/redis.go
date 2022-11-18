package receiver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisReceiver struct {
	*eventRouter
	conn *redis.Client
}

func NewRedisReceiver(connectOpts *redis.Options, opts ...ReceiverOption) (*RedisReceiver, error) {
	conn := redis.NewClient(connectOpts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := conn.Ping(ctx)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	router := newEventRouter(opts...)

	return &RedisReceiver{
		eventRouter: router,
		conn:        conn,
	}, nil
}

// Receive is a blocking function to receive events and pass them to handlers
func (r *RedisReceiver) Receive(ctx context.Context) (func(), error) {
	r.log.Debug().Msg("starting receive")
	pubsub := r.conn.PSubscribe(ctx, "discord.*")
	ch := pubsub.Channel()
	r.log.Debug().Str("pubsub", pubsub.String()).Msg("subscribed")

	stop := make(chan bool)

	func() {
		r.log.Debug().Msg("Listening for messages")
		for {
			select {
			case msg := <-ch:
				r.log.Debug().Str("channel", msg.Channel).Msg("received message")
				if err := r.Route(msg.Channel, json.RawMessage(msg.Payload)); err != nil {
					r.log.Warn().Err(err).Str("event", msg.Channel).Msg("failed to route event")
				}
			case <-stop:
				return
			}
		}
	}()

	return func() {
		stop <- true
		pubsub.Close()
	}, nil
}
