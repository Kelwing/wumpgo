package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisDispatcher struct {
	conn *redis.Client
}

func NewRedisDispatcher(connectOpts *redis.Options) (*RedisDispatcher, error) {
	rdb := redis.NewClient(connectOpts)

	return &RedisDispatcher{conn: rdb}, nil
}

func (d *RedisDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	log.Debug().Msgf("Dispatching event %s to Redis", eventName)
	cmd := d.conn.Publish(context.Background(), eventName, []byte(data))
	_, err := cmd.Result()
	return err
}
