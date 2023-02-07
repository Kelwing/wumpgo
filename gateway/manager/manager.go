package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/shard"
)

type ShardLocker struct {
	sync.Mutex
}

type ShardCluster struct {
	shards       []*shard.Shard
	shardCount   int
	concurrency  int
	clusterID    int
	shardOptions []shard.ShardOption
	token        string
	log          zerolog.Logger
}

func New(token string, o ...ManagerOption) *ShardCluster {
	m := &ShardCluster{
		shards:       make([]*shard.Shard, 0),
		shardCount:   1,
		concurrency:  1,
		clusterID:    0,
		shardOptions: []shard.ShardOption{},
		token:        token,
	}

	for _, opt := range o {
		opt(m)
	}

	return m
}

func (m *ShardCluster) createShards() error {
	if m.shardCount%m.concurrency != 0 {
		return fmt.Errorf("shard count must be a multiple of %d", m.concurrency)
	}

	shardsPerCluster := m.shardCount / m.concurrency
	m.shards = make([]*shard.Shard, shardsPerCluster)
	startID := m.concurrency * m.clusterID

	locker := &ShardLocker{}

	for i := 0; i < shardsPerCluster; i++ {
		staticOptions := []shard.ShardOption{
			shard.WithShardInfo(i+startID, m.shardCount),
			shard.WithIdentifyLock(locker),
		}
		m.shards[i] = shard.New(m.token,
			append(m.shardOptions, staticOptions...)...,
		)
	}

	return nil
}

func (m *ShardCluster) Run(ctx context.Context) error {
	err := m.createShards()
	if err != nil {
		return err
	}

	for _, s := range m.shards {
		go func(s *shard.Shard) {
			m.log.Info().Object("shard", s).Msg("starting shard")
			if err := s.Run(); err != nil {
				m.log.Error().Err(err).Object("shard", s).Msg("shard error")
			}
		}(s)
	}

	<-ctx.Done()

	for _, s := range m.shards {
		s.Close()
	}

	return nil
}
