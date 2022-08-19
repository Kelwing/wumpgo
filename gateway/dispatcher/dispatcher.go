package dispatcher

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type Dispatcher interface {
	Dispatch(event string, data json.RawMessage) error
}

type ErrUnknownDispatcher struct{}

func (e *ErrUnknownDispatcher) Error() string {
	return "Unknown Dispatcher"
}

func FromString(d string) (Dispatcher, error) {
	switch d {
	case "noop":
		return NewNOOPDispatcher(), nil
	case "nats":
		return NewNATSDispatcher(viper.GetString("dispatcher.addr"))
	case "redis":
		return NewRedisDispatcher(viper.GetString("dispatcher.addr"))
	case "nsq":
		return NewNSQDispatcher(viper.GetString("dispatcher.addr"))
	default:
		return nil, &ErrUnknownDispatcher{}
	}
}
