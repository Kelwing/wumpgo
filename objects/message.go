package objects

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=MessageType,MessageActivityType,MessageFlag,MessageStickerFormat -output message_string.go

type MessageType uint

const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeUserPremiumGuildSubscription
	MessageTypeUserPremiumGuildSubscriptionTier1
	MessageTypeUserPremiumGuildSubscriptionTier2
	MessageTypeUserPremiumGuildSubscriptionTier3
	MessageTypeChannelFollowAdd
	_
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	MessageTypeGuildDiscoveryGracePeriodInitialWarning
	MessageTypeGuildDiscoveryGracePeriodFinalWarning
	_
	MessageTypeReply
	MessageTypeApplicationCommand
)

type MessageActivityType uint

const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	_
	MessageActivityTypeJoinRequest
)

type MessageFlag uint

const (
	MsgFlagCrossposted MessageFlag = 1 << iota
	MsgFlagIsCrosspost
	MsgFlagSupressEmbeds
	MsgFlagSourceMessageDeleted
	MsgFlagUrgent
	_
	MsgFlagEphemeral
	MsgFlagLoading
)

type MessageStickerFormat uint

const (
	PNGStickerFormat MessageStickerFormat = iota + 1
	APNGStickerFormat
	LottieStickerFormat
)

type Message struct {
	ID                Snowflake           `json:"id"`
	ChannelID         Snowflake           `json:"channel_id"`
	GuildID           Snowflake           `json:"guild_id,omitempty"`
	Author            *User               `json:"author"`
	Member            *GuildMember        `json:"member,omitempty"`
	Content           string              `json:"content"`
	Timestamp         Time                `json:"timestamp"`
	EditedTimestamp   Time                `json:"edited_timestamp"`
	TTS               bool                `json:"tts"`
	MentionEveryone   bool                `json:"mention_everyone"`
	Mentions          []*User             `json:"mentions,omitempty"`
	MentionRoles      []Snowflake         `json:"mention_roles,omitempty"`
	MentionChannels   []*ChannelMention   `json:"mention_channels,omitempty"`
	Attachments       []*Attachment       `json:"attachments,omitempty"`
	Embeds            []*Embed            `json:"embeds"`
	Reactions         []*Reaction         `json:"reactions,omitempty"`
	Nonce             interface{}         `json:"nonce,omitempty"`
	Pinned            bool                `json:"pinned"`
	WebhookID         Snowflake           `json:"webhook_id,omitempty"`
	Type              MessageType         `json:"type"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *MessageApplication `json:"application,omitempty"`
	ApplicationID     Snowflake           `json:"application_id"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Flags             MessageFlag         `json:"flags,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
	Interaction       *MessageInteraction `json:"interaction,omitempty"`
	Thread            *Channel            `json:"thread,omitempty"`
	Components        []*Component        `json:"components,omitempty"`
	StickerItems      []*StickerItem      `json:"sticker_items"`
	Stickers          []*Sticker          `json:"stickers,omitempty"`
}

func (m *Message) URL() string {
	guild := "@me"
	if m.GuildID != Snowflake(0) {
		guild = strconv.FormatUint(uint64(m.GuildID), 10)
	}
	return fmt.Sprintf("https://discord.com/channels/%s/%d/%d", guild, m.ChannelID, m.ID)
}

type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User *User           `json:"user"`
}

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id,omitempty"`
}

type MessageReference struct {
	MessageID       Snowflake `json:"message_id,omitempty"`
	ChannelID       Snowflake `json:"channel_id,omitempty"`
	GuildID         Snowflake `json:"guild_id,omitempty"`
	FailIfNotExists *bool     `json:"fail_if_not_exists,omitempty"`
}

type MessageApplication struct {
	ID          Snowflake `json:"id"`
	CoverImage  string    `json:"cover_image,omitempty"`
	Description string    `json:"description"`
	Icon        string    `json:"icon,omitempty"`
	Name        string    `json:"name"`
}

type Attachment struct {
	ID          Snowflake `json:"id"`
	Filename    string    `json:"filename,omitempty"`
	Description string    `json:"description,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
	Size        int       `json:"size,omitempty"`
	URL         string    `json:"url,omitempty"`
	ProxyURL    string    `json:"proxy_url,omitempty"`
	Height      int       `json:"height,omitempty"`
	Width       int       `json:"width,omitempty"`
	Ephemeral   bool      `json:"ephemeral,omitempty"`
}

type ChannelMention struct {
	ID      Snowflake   `json:"id"`
	GuildID Snowflake   `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}
