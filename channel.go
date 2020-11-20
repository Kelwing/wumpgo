package objects

type ChannelType uint

const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildNews
	ChannelTypeGuildStore
)

type PermissionOverwrite struct {
	ID    Snowflake `json:"id"`
	Type  uint      `json:"type"`
	Allow string    `json:"allow"`
	Deny  string    `json:"deny"`
}

type Channel struct {
	ID                   Snowflake             `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              Snowflake             `json:"guild_id,omitempty"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Topic                string                `json:"topic,omitempty"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
	LastMessageID        Snowflake             `json:"last_message_id,omitempty"`
	Bitrate              uint                  `json:"bitrate,omitempty"`
	UserLimit            uint                  `json:"user_limit,omitempty"`
	RateLimitPerUser     uint                  `json:"rate_limit_per_user,omitempty"`
	Recipients           []*User               `json:"recipient,omitempty"`
	Icon                 string                `json:"icon,omitempty"`
	OwnerID              Snowflake             `json:"owner_id,omitempty"`
	ApplicationID        Snowflake             `json:"application_id,omitempty"`
	ParentID             Snowflake             `json:"parent_id,omitempty"`
	LastPinTimestamp     Time                  `json:"last_pin_timestamp,omitempty"`
}

type AllowedMentions struct {
	Parse       []string    `json:"parse"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Users       []Snowflake `json:"users,omitempty"`
	RepliedUser bool        `json:"replied_user,omitempty"`
}
