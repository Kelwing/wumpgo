package restcache

import (
	"context"

	"github.com/dgrijalva/lfu-go"
	"wumpgo.dev/wumpgo/rest"
)

type MemoryCache struct {
	c *lfu.Cache
}

type MemoryCacheOption func(c *MemoryCache)

func WithUpperBound(b int) MemoryCacheOption {
	return func(c *MemoryCache) {
		c.c.UpperBound = b
	}
}

func WithLowerBount(b int) MemoryCacheOption {
	return func(c *MemoryCache) {
		c.c.LowerBound = b
	}
}

func NewMemoryCache(opts ...MemoryCacheOption) *MemoryCache {
	c := &MemoryCache{
		c: lfu.New(),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *MemoryCache) Get(ctx context.Context, key string) (*rest.DiscordResponse, error) {
	i := c.c.Get(key)
	resp, ok := i.(*rest.DiscordResponse)
	if !ok {
		return nil, nil
	}

	return resp, nil
}

func (c *MemoryCache) Put(ctx context.Context, key string, value *rest.DiscordResponse) error {
	c.c.Set(key, value)
	return nil
}
