package shard

import (
	"fmt"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/objects"
)

type ShardOpt interface {
	apply(s *Shard)
}

type withGatewayURL struct {
	u string
}

func (o *withGatewayURL) apply(s *Shard) {
	s.gateway_url = fmt.Sprintf(GatewayAddressFmt, o.u, GatewayVersion, GatewayEncoding)
}

func WithGatewayURL(u string) ShardOpt {
	return &withGatewayURL{u: u}
}

type withDispatcher struct {
	d dispatcher.Dispatcher
}

func (o *withDispatcher) apply(s *Shard) {
	s.dispatcher = o.d
}

func WithDispatcher(d dispatcher.Dispatcher) ShardOpt {
	return &withDispatcher{d: d}
}

type withShardInfo struct {
	id    int
	count int
}

func (o *withShardInfo) apply(s *Shard) {
	s.identify.Shard = []int{o.id, o.count}
}

func WithShardInfo(id, count int) ShardOpt {
	return &withShardInfo{id: id, count: count}
}

type withIntents struct {
	intents objects.Intent
}

func (o *withIntents) apply(s *Shard) {
	s.identify.Intents = o.intents
}

func WithIntents(i objects.Intent) ShardOpt {
	return &withIntents{intents: i}
}

type withLogger struct {
	logger zerolog.Logger
}

func (o *withLogger) apply(s *Shard) {
	s.logger = o.logger
}

func WithLogger(l zerolog.Logger) ShardOpt {
	return &withLogger{logger: l}
}

type withInitialPresence struct {
	presence objects.UpdatePresence
}

func (o *withInitialPresence) apply(s *Shard) {
	s.identify.Presence = o.presence
}

func WithInitialPresence(p objects.UpdatePresence) ShardOpt {
	return &withInitialPresence{presence: p}
}
