package dispatcher

import (
	"encoding/json"
)

type Dispatcher interface {
	Dispatch(event string, data json.RawMessage) error
}
