package objects

type PrivacyLevel int

const (
	PrivacyLevel_PUBLIC = iota + 1
	PrivacyLevel_GUILD_ONLY
)

type StageInstance struct {
	ID                   Snowflake    `json:"id"`
	GuildID              Snowflake    `json:"guild_id"`
	ChannelID            Snowflake    `json:"channel_id"`
	Topic                string       `json:"topic"`
	PrivacyLevel         PrivacyLevel `json:"privacy_level"`
	DiscoverableDisabled bool         `json:"discoverable_disabled"`
}
