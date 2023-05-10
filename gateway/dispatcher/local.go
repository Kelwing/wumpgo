package dispatcher

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"wumpgo.dev/wumpgo/gateway/receiver"
)

var _ Dispatcher = (*LocalDispatcher)(nil)

type LocalDispatcher struct {
	receiver receiver.Receiver
	logger   *zerolog.Logger
}

func NewLocalDispatcher(receiver receiver.Receiver, opts ...DispatcherOption) *LocalDispatcher {
	logger := zerolog.Nop()

	l := &LocalDispatcher{
		receiver: receiver,
		logger:   &logger,
	}

	for _, o := range opts {
		o(l)
	}

	return l
}

func (l *LocalDispatcher) Dispatch(event string, data json.RawMessage) error {
	return l.receiver.Route(event, data)
}

func (l *LocalDispatcher) SetLogger(logger *zerolog.Logger) {
	l.logger = logger
}
