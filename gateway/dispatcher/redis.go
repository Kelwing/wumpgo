package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ Dispatcher = (*RedisDispatcher)(nil)

type RedisDispatcher struct {
	conn   *redis.Client
	logger *zerolog.Logger
}

func NewRedisDispatcher(connectOpts *redis.Options, opts ...DispatcherOption) (*RedisDispatcher, error) {
	logger := zerolog.Nop()
	rdb := redis.NewClient(connectOpts)

	d := &RedisDispatcher{conn: rdb, logger: &logger}

	for _, o := range opts {
		o(d)
	}

	return d, nil
}

func (d *RedisDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	log.Debug().Msgf("Dispatching event %s to Redis", eventName)
	cmd := d.conn.Publish(context.Background(), eventName, []byte(data))
	_, err := cmd.Result()
	return err
}

func (d *RedisDispatcher) setLogger(logger *zerolog.Logger) {
	d.logger = logger
}
