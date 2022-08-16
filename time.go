package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var _ json.Marshaler = (*Time)(nil)
var _ json.Unmarshaler = (*Time)(nil)
var _ Mentionable = (*Time)(nil)

type TimestampStyle string

const (
	StyleShortTime     TimestampStyle = "t"
	StyleLongTime      TimestampStyle = "T"
	StyleShortDate     TimestampStyle = "d"
	StyleLongDate      TimestampStyle = "D"
	StyleShortDateTime TimestampStyle = "f"
	StyleLongDateTime  TimestampStyle = "F"
	StyleRelative      TimestampStyle = "R"
)

func NewTime(t time.Time) *Time {
	return &Time{t}
}

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(t.Time.Format(time.RFC3339))), nil
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

func (t *Time) Format(style TimestampStyle) string {
	return fmt.Sprintf("<t:%d:%s>", t.Unix(), style)
}

func (t *Time) String() string {
	return t.ShortDateTime()
}

func (t *Time) Mention() string {
	return t.LongDateTime()
}

func (t *Time) ShortTime() string {
	return t.Format(StyleShortTime)
}

func (t *Time) LongTime() string {
	return t.Format(StyleLongTime)
}

func (t *Time) ShortDate() string {
	return t.Format(StyleShortDate)
}

func (t *Time) LongDate() string {
	return t.Format(StyleLongDate)
}

func (t *Time) ShortDateTime() string {
	return t.Format(StyleShortDateTime)
}

func (t *Time) LongDateTime() string {
	return t.Format(StyleLongDateTime)
}

func (t *Time) Relative() string {
	return t.Format(StyleRelative)
}
