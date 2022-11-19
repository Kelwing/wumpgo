package dispatcher

import (
	"encoding/json"

	"wumpgo.dev/wumpgo/gateway/receiver"
)

type LocalDispatcher struct {
	receiver *receiver.LocalReceiver
}

func NewLocalDispatcher(receiver *receiver.LocalReceiver) *LocalDispatcher {
	return &LocalDispatcher{
		receiver: receiver,
	}
}

func (l *LocalDispatcher) Dispatch(event string, data json.RawMessage) error {
	return l.receiver.Route(event, data)
}
