package dispatcher

import (
	"encoding/json"

	"wumpgo.dev/wumpgo/gateway/receiver"
)

type LocalDispatcher struct {
	receiver receiver.Receiver
}

func NewLocalDispatcher(receiver receiver.Receiver) *LocalDispatcher {
	return &LocalDispatcher{
		receiver: receiver,
	}
}

func (l *LocalDispatcher) Dispatch(event string, data json.RawMessage) error {
	return l.receiver.Route(event, data)
}
