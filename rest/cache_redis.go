package rest

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cache = (*RedisCache)(nil)

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

type RedisCache struct {
	rdb *redis.Client
}

func (r *RedisCache) Get(ctx context.Context, key string) (*DiscordResponse, error) {
	resp := r.rdb.Get(ctx, key)
	if resp.Err() != nil {
		return nil, &CacheError{message: resp.Err().Error()}
	}

	b, err := resp.Bytes()
	if err != nil {
		return nil, &CacheError{message: err.Error()}
	}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	var discordResp DiscordResponse

	if err := dec.Decode(&discordResp); err != nil {
		return nil, &CacheError{message: err.Error()}
	}

	return &discordResp, nil
}

func (r *RedisCache) Put(ctx context.Context, key string, value *DiscordResponse) error {
	buf := &bytes.Buffer{}

	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		return &CacheError{message: err.Error()}
	}

	resp := r.rdb.Set(ctx, key, buf.String(), time.Minute)
	if resp.Err() != nil {
		return &CacheError{message: resp.Err().Error()}
	}

	return nil
}
