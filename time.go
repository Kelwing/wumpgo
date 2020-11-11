package objects

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(t.String())), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte(`""`)) {
		return nil
	}
	var ts time.Time
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}
	t.Time = ts
	return nil
}

func (t *Time) String() string {
	return t.Time.Format(time.RFC3339)
}
