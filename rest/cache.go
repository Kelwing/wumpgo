package rest

import (
	"context"
	"errors"
	"sync"
	"time"
)

type CacheError struct {
	message string
}

func (c *CacheError) Error() string {
	return c.message
}

type Cache interface {
	Get(ctx context.Context, key string) (*DiscordResponse, error)
	Put(ctx context.Context, key string, value *DiscordResponse) error
}

var _ Cache = (*MemoryCache)(nil)

// Reference in-memory implementation
type MemoryCache struct {
	sync.RWMutex
	cacheMap  map[string]*CacheValue
	config    *MemoryCacheConfig
	cleanChan chan string
}

type CacheValue struct {
	resp *DiscordResponse
	// insertion time.Time
	timer *time.Timer
}

type MemoryCacheConfig struct {
	retention time.Duration
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		cacheMap:  make(map[string]*CacheValue),
		cleanChan: make(chan string),
	}

	go cache.cleaner()

	return cache
}

func (m *MemoryCache) Get(ctx context.Context, key string) (*DiscordResponse, error) {
	m.RLock()
	defer m.RUnlock()
	data, ok := m.cacheMap[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return data.resp, nil
}

func (m *MemoryCache) Put(ctx context.Context, key string, resp *DiscordResponse) error {
	m.Lock()
	defer m.Unlock()

	m.cacheMap[key] = &CacheValue{
		resp: resp,
		timer: time.AfterFunc(m.config.retention, func() {
			m.cleanChan <- key
		}),
	}
	return nil
}

func (m *MemoryCache) cleaner() {
	for key := range m.cleanChan {
		m.Lock()
		delete(m.cacheMap, key)
		m.Unlock()
	}
}
