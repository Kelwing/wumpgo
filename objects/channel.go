package objects

//go:generate stringer -type=ChannelType -trimprefix=ChannelType -output channel_string.go

var _ Mentionable = (*Channel)(nil)
var _ SnowflakeObject = (*Channel)(nil)
var _ SnowflakeObject = (*ThreadMember)(nil)

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
	DiscordBaseObject
	Type                       ChannelType           `json:"type"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      string                `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              Snowflake             `json:"last_message_id,omitempty"`
	Bitrate                    uint                  `json:"bitrate,omitempty"`
	UserLimit                  uint                  `json:"user_limit,omitempty"`
	RateLimitPerUser           uint                  `json:"rate_limit_per_user,omitempty"`
	Recipients                 []*User               `json:"recipient,omitempty"`
	Icon                       string                `json:"icon,omitempty"`
	OwnerID                    Snowflake             `json:"owner_id,omitempty"`
	ApplicationID              Snowflake             `json:"application_id,omitempty"`
	ParentID                   Snowflake             `json:"parent_id,omitempty"`
	LastPinTimestamp           Time                  `json:"last_pin_timestamp,omitempty"`
	RtcRegion                  *string               `json:"rtc_region,omitempty"`
	VideoQualityMode           *int                  `json:"video_quality_mode,omitempty"`
	MessageCount               *int                  `json:"message_count,omitempty"`
	MemberCount                *int                  `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration *int                  `json:"default_auto_archive_duration,omitempty"`
	Permissions                *string               `json:"permissions,omitempty"`
}

func (c *Channel) Mention() string {
	return "<#" + c.GetID().String() + ">"
}

type AllowedMentions struct {
	Parse       []string    `json:"parse"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Users       []Snowflake `json:"users,omitempty"`
	RepliedUser bool        `json:"replied_user,omitempty"`
}

type FollowedChannel struct {
	ChannelID Snowflake `json:"channel_id"`
	WebhookID Snowflake `json:"webhook_id"`
}

type ThreadMember struct {
	DiscordBaseObject
	UserID   Snowflake `json:"user_id"`
	JoinedAt Time      `json:"join_timestamp"`
	Flags    uint      `json:"flags"`
}

type ThreadMetadata struct {
	Archived            bool `json:"archived"`
	AutoArchiveDuration int  `json:"auto_archive_duration"`
	ArchivedTimestamp   Time `json:"archived_timestamp"`
	Locked              bool `json:"locked"`
	Invitable           bool `json:"invitable"`
}
