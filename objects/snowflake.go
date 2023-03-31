package objects

import (
	"strconv"
	"time"

	"wumpgo.dev/snowflake"
)

type Snowflake snowflake.Snowflake

var _ SnowflakeObject = (*Snowflake)(nil)

const (
	DiscordEpoch = 1420070400000
)

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

func SnowflakeFromString(s string) (Snowflake, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return Snowflake(i), nil
}
