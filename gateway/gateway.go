package gateway

import (
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"wumpgo.dev/wumpgo/gateway/dispatcher"
	"wumpgo.dev/wumpgo/gateway/manager"
	"wumpgo.dev/wumpgo/gateway/receiver"
	"wumpgo.dev/wumpgo/gateway/shard"
)

func NewShard(token string, opts ...shard.ShardOption) *shard.Shard {
	return shard.New(token, opts...)
}

func NewShardManager(token string, opts ...manager.ManagerOption) *manager.ShardCluster {
	return manager.New(token, opts...)
}

func NewLocalDispatcher(receiver receiver.Receiver, opts ...dispatcher.DispatcherOption) dispatcher.Dispatcher {
	return dispatcher.NewLocalDispatcher(receiver)
}

func NewNATSDispatcher(url string, natsOpts []nats.Option, opts ...dispatcher.DispatcherOption) (dispatcher.Dispatcher, error) {
	return dispatcher.NewNATSDispatcher(url, natsOpts, opts...)
}

func NewNOOPDispatcher(opts ...dispatcher.DispatcherOption) dispatcher.Dispatcher {
	return dispatcher.NewNOOPDispatcher()
}

func NewRedisDispatcher(connectOpts *redis.Options, opts ...dispatcher.DispatcherOption) (dispatcher.Dispatcher, error) {
	return dispatcher.NewRedisDispatcher(connectOpts)
}

func NewLocalReceiver(opts ...receiver.ReceiverOption) receiver.Receiver {
	return receiver.NewLocalReceiver(opts...)
}

func NewNATSReceiver(url string, natsOptions []nats.Option, opts ...receiver.ReceiverOption) (receiver.Receiver, error) {
	return receiver.NewNATSReceiver(url, natsOptions, opts...)
}

func NewRedisReceiver(connectOpts *redis.Options, opts ...receiver.ReceiverOption) (receiver.Receiver, error) {
	return receiver.NewRedisReceiver(connectOpts, opts...)
}
