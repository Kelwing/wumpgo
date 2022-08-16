package objects

//go:generate stringer -type=PrivacyLevel -trimprefix=PrivacyLevel -output stage_string.go

var _ SnowflakeObject = (*StageInstance)(nil)

type PrivacyLevel int

const (
	PrivacyLevelPublic PrivacyLevel = iota + 1
	PrivacyLevelGuildOnly
)

type StageInstance struct {
	DiscordBaseObject
	GuildID              Snowflake    `json:"guild_id"`
	ChannelID            Snowflake    `json:"channel_id"`
	Topic                string       `json:"topic"`
	PrivacyLevel         PrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled bool         `json:"discoverable_disabled"`
}
