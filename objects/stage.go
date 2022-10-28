package objects

//go:generate stringer -type=PrivacyLevel -trimprefix=PrivacyLevel -output stage_string.go

type PrivacyLevel int

const (
	PrivacyLevelPublic PrivacyLevel = iota + 1
	PrivacyLevelGuildOnly
)

type StageInstance struct {
	ID                    Snowflake    `json:"id"`
	GuildID               Snowflake    `json:"guild_id"`
	ChannelID             Snowflake    `json:"channel_id"`
	Topic                 string       `json:"topic"`
	PrivacyLevel          PrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled  bool         `json:"discoverable_disabled"`
	GuildScheduledEventID *Snowflake   `json:"guild_scheduled_event_id,omitempty"`
}
