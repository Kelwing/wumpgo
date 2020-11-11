package objects

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	DiscordEpoch = 1420070400000
)

type Snowflake int64

func (s *Snowflake) UnmarshalJSON(bytes []byte) error {
	var snowflake string
	err := json.Unmarshal(bytes, &snowflake)
	if err != nil {
		return err
	}

	snowInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return err
	}

	*s = Snowflake(snowInt)

	return nil
}

func (s Snowflake) MarshalJSON() ([]byte, error) {
	snowString := strconv.Itoa(int(s))
	return json.Marshal(snowString)
}

// CreatedAt returns a time.Time representing the time a Snowflake was created
func (s Snowflake) CreatedAt() time.Time {
	timestampMs := (int64(s) >> 22) + DiscordEpoch
	return time.Unix(0, timestampMs*int64(time.Millisecond))
}
