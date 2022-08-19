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

func NewRedisDispatcher(url string) (*RedisDispatcher, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: url,
		DB:   0,
	})

	return &RedisDispatcher{conn: rdb}, nil
}

func (d *RedisDispatcher) Dispatch(event string, data json.RawMessage) error {
	eventName := fmt.Sprintf("discord.%s", strings.ToLower(event))
	log.Debug().Msgf("Dispatching event %s to Redis", eventName)
	return d.conn.Publish(context.Background(), eventName, data).Err()
}
