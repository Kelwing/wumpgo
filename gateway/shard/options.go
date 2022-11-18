package shard

import (
	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/objects"
)

type ShardOption func(*Shard)

func WithGatewayURL(u string) ShardOption {
	return func(s *Shard) {
		s.gateway_url = u
	}
}

func WithDispatcher(d dispatcher.Dispatcher) ShardOption {
	return func(s *Shard) {
		s.dispatcher = d
	}
}

func WithShardInfo(id, count int) ShardOption {
	return func(s *Shard) {
		s.identify.Shard = []int{id, count}
	}
}

func WithIntents(i objects.Intent) ShardOption {
	return func(s *Shard) {
		s.identify.Intents = i
	}
}

func WithLogger(l zerolog.Logger) ShardOption {
	return func(s *Shard) {
		s.logger = l
	}
}

func WithInitialPresence(p objects.UpdatePresence) ShardOption {
	return func(s *Shard) {
		s.identify.Presence = p
	}
}
