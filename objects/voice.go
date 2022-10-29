package objects

type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	VIP        bool   `json:"vip"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

type VoiceState struct {
	GuildID                 Snowflake    `json:"guild_id,omitempty"`
	ChannelID               Snowflake    `json:"channel_id,omitempty"`
	UserID                  Snowflake    `json:"user_id"`
	Member                  *GuildMember `json:"guild_member,omitempty"`
	SessionID               string       `json:"session_id"`
	Deaf                    bool         `json:"deaf"`
	Mute                    bool         `json:"mute"`
	SelfDeaf                bool         `json:"self_deaf"`
	SelfMute                bool         `json:"self_mute"`
	SelfStream              bool         `json:"self_stream,omitempty"`
	SelfVideo               bool         `json:"self_video"`
	Suppress                bool         `json:"suppress"`
	RequestToSpeakTimestamp *Time        `json:"request_to_speak_timestamp,omitempty"`
}
