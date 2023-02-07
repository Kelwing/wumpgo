package manager

import "wumpgo.dev/wumpgo/gateway/shard"

type ManagerOption func(m *ShardCluster)

func WithShardCount(c int) ManagerOption {
	return func(m *ShardCluster) {
		m.shardCount = c
	}
}

func WithConcurrency(c int) ManagerOption {
	return func(m *ShardCluster) {
		m.concurrency = c
	}
}

func WithClusterID(i int) ManagerOption {
	return func(m *ShardCluster) {
		m.clusterID = i
	}
}

func WithShardOptions(o ...shard.ShardOption) ManagerOption {
	return func(m *ShardCluster) {
		m.shardOptions = o
	}
}
