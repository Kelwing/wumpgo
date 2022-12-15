package restcache

import (
	"context"
	"time"

	rcache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"wumpgo.dev/wumpgo/rest"
)

var _ rest.Cache = (*RedisCache)(nil)

type RedisCacheOption func(r *RedisCache)

func WithTTL(ttl time.Duration) RedisCacheOption {
	return func(r *RedisCache) {
		r.ttl = ttl
	}
}

func WithRedisClient(c *redis.Client) RedisCacheOption {
	return func(r *RedisCache) {
		r.opts.Redis = c
	}
}

func WithRedisRing(c *redis.Ring) RedisCacheOption {
	return func(r *RedisCache) {
		r.opts.Redis = c
	}
}

func WithLocalCache(lc rcache.LocalCache) RedisCacheOption {
	return func(r *RedisCache) {
		r.opts.LocalCache = lc
	}
}

type RedisCache struct {
	opts  *rcache.Options
	cache *rcache.Cache
	ttl   time.Duration
}

func NewRedisCache(opts ...RedisCacheOption) *RedisCache {
	c := &RedisCache{
		ttl: time.Minute,
		opts: &rcache.Options{
			Redis: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			}),
		},
	}

	for _, o := range opts {
		o(c)
	}

	c.cache = rcache.New(c.opts)

	return c
}

func (r *RedisCache) Get(ctx context.Context, key string) (*rest.DiscordResponse, error) {
	var obj *rest.DiscordResponse

	err := r.cache.Get(ctx, key, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *RedisCache) Put(ctx context.Context, key string, value *rest.DiscordResponse) error {
	return r.cache.Set(&rcache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   r.ttl,
	})
}
