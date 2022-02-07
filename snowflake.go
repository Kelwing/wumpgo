package objects

import (
	"encoding/json"
	"strconv"
	"time"
)

var _ SnowflakeObject = (*Snowflake)(nil)

const (
	DiscordEpoch = 1420070400000
)

type Snowflake uint64

func (s *Snowflake) UnmarshalJSON(bytes []byte) error {
	var snowflake string
	err := json.Unmarshal(bytes, &snowflake)
	if err != nil {
		return err
	}

	if snowflake == "" || snowflake == "null" {
		*s = 0
		return nil
	}

	snowInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return err
	}

	*s = Snowflake(snowInt)

	return nil
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// CreatedAt returns a time.Time representing the time a Snowflake was created
func (s Snowflake) CreatedAt() Time {
	timestampMs := (int64(s) >> 22) + DiscordEpoch
	return Time{time.Unix(0, timestampMs*int64(time.Millisecond))}
}

func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s Snowflake) GetID() Snowflake {
	return s
}
