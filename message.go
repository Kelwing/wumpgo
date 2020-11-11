package objects

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
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
)

type MessageActivityType uint

const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest = 5
)

type MessageFlag uint

const (
	MsgFlagCrossposted MessageFlag = 1 << iota
	MsgFlagIsCrosspost
	MsgFlagSupressEmbeds
	MsgFlagSourceMessageDeleted
	MsgFlagUrgent
)

type Message struct {
	ID               Snowflake           `json:"id"`
	ChannelID        Snowflake           `json:"channel_id"`
	GuildID          Snowflake           `json:"guild_id,omitempty"`
	Author           *User               `json:"author"`
	Member           *GuildMember        `json:"member,omitempty"`
	Content          string              `json:"content"`
	Timestamp        Time                `json:"timestamp"`
	EditedTimestamp  Time                `json:"edited_timestamp"`
	TTS              bool                `json:"tts"`
	MentionEveryone  bool                `json:"mention_everyone"`
	Mentions         []*User             `json:"mentions,omitempty"`
	MentionRoles     []Snowflake         `json:"mention_roles,omitempty"`
	MentionChannels  []*ChannelMention   `json:"mention_channels,omitempty"`
	Attachments      []*Attachment       `json:"attachments,omitempty"`
	Embeds           []*Embed            `json:"embeds"`
	Reactions        []*Reaction         `json:"reactions,omitempty"`
	Nonce            interface{}         `json:"nonce,omitempty"`
	Pinned           bool                `json:"pinned"`
	WebhookID        Snowflake           `json:"webhook_id,omitempty"`
	Type             MessageType         `json:"type"`
	Activity         *MessageActivity    `json:"activity,omitempty"`
	Application      *MessageApplication `json:"application,omitempty"`
	MessageReference *MessageReference   `json:"message_reference,omitempty"`
	Flags            MessageFlag         `json:"flags,omitempty"`
}

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id,omitempty"`
}

type MessageReference struct {
	MessageID Snowflake `json:"message_id,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
}

type MessageApplication struct {
	ID          Snowflake `json:"id"`
	CoverImage  string    `json:"cover_image,omitempty"`
	Description string    `json:"description"`
	Icon        string    `json:"icon,omitempty"`
	Name        string    `json:"name"`
}

type Attachment struct {
	ID       Snowflake `json:"id"`
	Filename string    `json:"filename"`
	Size     int       `json:"size"`
	URL      string    `json:"url"`
	ProxyURL string    `json:"proxy_url"`
	Height   int       `json:"height,omitempty"`
	Width    int       `json:"width,omitempty"`
}

type ChannelMention struct {
	ID      Snowflake   `json:"id"`
	GuildID Snowflake   `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}
