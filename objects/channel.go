package objects

import (
	"wumpgo.dev/wumpgo/objects/permissions"
)

//go:generate stringer -type=ChannelType,ChannelFlag,ChannelForumSortOrder -trimprefix=Channel -output channel_string.go

var _ Mentionable = (*Channel)(nil)

type ChannelType uint

const (
	ChannelTypeGuildText ChannelType = iota
	ChannelTypeDM
	ChannelTypeGuildVoice
	ChannelTypeGroupDM
	ChannelTypeGuildCategory
	ChannelTypeGuildAnnouncement
	_
	_
	_
	_
	ChannelTypeAnnouncementThread
	ChannelTypePublicThread
	ChannelTypePrivateThread
	ChannelTypeGuildStageVoice
	ChannelTypeGuildDirectory
	ChannelTypeGuildForum
)

type PermissionOverwrite struct {
	ID    Snowflake `json:"id"`
	Type  uint      `json:"type"`
	Allow string    `json:"allow"`
	Deny  string    `json:"deny"`
}

type Channel struct {
	ID                            Snowflake                 `json:"id"`
	Type                          ChannelType               `json:"type"`
	GuildID                       Snowflake                 `json:"guild_id,omitempty"`
	Position                      int                       `json:"position,omitempty"`
	PermissionOverwrites          []PermissionOverwrite     `json:"permission_overwrites,omitempty"`
	Name                          string                    `json:"name,omitempty"`
	Topic                         string                    `json:"topic,omitempty"`
	NSFW                          bool                      `json:"nsfw,omitempty"`
	LastMessageID                 Snowflake                 `json:"last_message_id,omitempty"`
	Bitrate                       uint                      `json:"bitrate,omitempty"`
	UserLimit                     uint                      `json:"user_limit,omitempty"`
	RateLimitPerUser              uint                      `json:"rate_limit_per_user,omitempty"`
	Recipients                    []*User                   `json:"recipient,omitempty"`
	Icon                          string                    `json:"icon,omitempty"`
	OwnerID                       Snowflake                 `json:"owner_id,omitempty"`
	ApplicationID                 Snowflake                 `json:"application_id,omitempty"`
	ParentID                      Snowflake                 `json:"parent_id,omitempty"`
	LastPinTimestamp              Time                      `json:"last_pin_timestamp,omitempty"`
	RtcRegion                     *string                   `json:"rtc_region,omitempty"`
	VideoQualityMode              *int                      `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                      `json:"message_count,omitempty"`
	MemberCount                   *int                      `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata           `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember             `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int                      `json:"default_auto_archive_duration,omitempty"`
	Permissions                   permissions.PermissionBit `json:"permissions,omitempty"`
	Flags                         ChannelFlag               `json:"flags,omitempty"`
	TotalMessagesSent             uint                      `json:"total_message_sent,omitempty"`
	AvailableTags                 []*ForumTag               `json:"available_tags,omitempty"`
	AppliedTags                   []Snowflake               `json:"applied_tags,omitempty"`
	DefaultReactionEmoji          *DefaultReaction          `json:"default_reaction_emoji,omitempty"`
	DefaultThreadRateLimitPerUser int                       `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultSortOrder              *int                      `json:"default_sort_order,omitempty"`
}

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

type ForumThreadChannel struct {
	*Channel
	Message *Message
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
	ID       Snowflake `json:"id"`
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

type ChannelFlag uint64

const (
	ChannelFlagPinned ChannelFlag = 1 << (iota + 1)
	_
	_
	ChannelFlagRequireTag
)

type ForumTag struct {
	ID        Snowflake `json:"id"`
	Name      string    `json:"name"`
	Moderated bool      `json:"moderated"`
	EmojiID   Snowflake `json:"emoji_id"`
	EmojiName string    `json:"string,omitempty"`
}

type DefaultReaction struct {
	EmojiID   Snowflake `json:"emoji_id,omitempty"`
	EmojiName string    `json:"emoji_name,omitempty"`
}

type ChannelForumSortOrder uint64

const (
	ChannelForumSortOrderLatestActivity ChannelForumSortOrder = iota
	ChannelForumSortOrderCreationDate
)
