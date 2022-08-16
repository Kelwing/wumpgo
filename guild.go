package objects

import "github.com/Postcord/objects/permissions"

var _ SnowflakeObject = (*Guild)(nil)
var _ SnowflakeObject = (*UnavailableGuild)(nil)
var _ SnowflakeObject = (*GuildPreview)(nil)
var _ SnowflakeObject = (*WidgetUser)(nil)
var _ SnowflakeObject = (*GuildWidgetJSON)(nil)

//go:generate stringer -type=VerificationLevel,MessageNotificationsLevel,ExplicitContentFilterLevel,MFALevel,PremiumTierLevel,SystemChannelFlags,GuildNSFWLevel -output guild_string.go

type VerificationLevel uint

const (
	VerificationLevelNone VerificationLevel = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

type MessageNotificationsLevel uint

const (
	MessageNotificationsLevelAll MessageNotificationsLevel = iota
	MessageNotificationsLevelOnlyMentions
)

type ExplicitContentFilterLevel uint

const (
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota
	ExplicitContentFilterLevelMembersWithoutRoles
	ExplicitContentFilterLevelAllMembers
)

type GuildFeature string

const (
	GuildFeatureAnimatedIcon                  GuildFeature = "ANIMATED_ICON"
	GuildFeatureBanner                        GuildFeature = "BANNER"
	GuildFeatureCommerce                      GuildFeature = "COMMERCE"
	GuildFeatureCommunity                     GuildFeature = "COMMUNITY"
	GuildFeatureDiscoverable                  GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable                    GuildFeature = "FEATURABLE"
	GuildFeatureInviteSplash                  GuildFeature = "INVITE_SPLASH"
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	GuildFeatureMonetizationEnabled           GuildFeature = "MONETIZATION_ENABLED"
	GuildFeatureMoreStickers                  GuildFeature = "MORE_STICKERS"
	GuildFeatureNews                          GuildFeature = "NEWS"
	GuildFeaturePartnered                     GuildFeature = "PARTNERED"
	GuildFeaturePreviewEnabled                GuildFeature = "PREVIEW_ENABLED"
	GuildFeaturePrivateThreads                GuildFeature = "PRIVATE_THREADS"
	GuildFeatureRoleIcons                     GuildFeature = "ROLE_ICONS"
	GuildFeatureSevenDayThreadArchive         GuildFeature = "SEVEN_DAY_THREAD_ARCHIVE"
	GuildFeatureThreeDayThreadArchive         GuildFeature = "THREE_DAY_THREAD_ARCHIVE"
	GuildFeatureTicketedEventsEnabled         GuildFeature = "TICKETED_EVENTS_ENABLED"
	GuildFeatureVanityURL                     GuildFeature = "VANITY_URL"
	GuildFeatureVerified                      GuildFeature = "VERIFIED"
	GuildFeatureVIPRegions                    GuildFeature = "VIP_REGIONS"
	GuildFeatureWelcomeScreenEnabled          GuildFeature = "WELCOME_SCREEN_ENABLED"
)

type MFALevel uint

const (
	MFALevelNone MFALevel = iota
	MFALevelElevated
)

type PremiumTierLevel uint

const (
	PremiumTierNone PremiumTierLevel = iota
	PremiumTier1
	PremiumTier2
	PremiumTier3
)

type SystemChannelFlags uint

const (
	FlagSupressJoinNotifications SystemChannelFlags = 1 << iota
	FlagSupressPremiumSubscriptions
)

type GuildNSFWLevel int

const (
	GuildNSFWLevelDefault GuildNSFWLevel = iota
	GuildNSFWLevelExplicit
	GuildNSFWLevelSafe
	GuildNSFWLevelAgeRestricted
)

type Guild struct {
	DiscordBaseObject
	Name                        string                     `json:"name"`
	Icon                        string                     `json:"icon,omitempty"`
	IconHash                    string                     `json:"icon_hash,omitempty"`
	Splash                      string                     `json:"splash,omitempty"`
	DiscoverySplash             string                     `json:"discovery_splash,omitempty"`
	Owner                       bool                       `json:"owner,omitempty"`
	OwnerID                     Snowflake                  `json:"owner_id"`
	Permissions                 permissions.PermissionBit  `json:"permissions,omitempty"`
	Region                      string                     `json:"region"`
	AFKChannelID                Snowflake                  `json:"afk_channel_id,omitempty"`
	AFKTimeout                  int                        `json:"afk_timeout,omitempty"`
	WidgetEnabled               bool                       `json:"widget_enabled,omitempty"`
	WidgetChannelID             Snowflake                  `json:"widget_channel_id,omitempty"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications MessageNotificationsLevel  `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Roles                       []*Role                    `json:"roles"`
	Emojis                      []*Emoji                   `json:"emojis"`
	Features                    []GuildFeature             `json:"features"`
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               Snowflake                  `json:"application_id,omitempty"`
	SystemChannelID             Snowflake                  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          SystemChannelFlags         `json:"system_channel_flags"`
	RulesChannelID              Snowflake                  `json:"rules_channel_id,omitempty"`
	JoinedAt                    Time                       `json:"joined_at,omitempty"`
	Large                       bool                       `json:"large,omitempty"`
	Unavailable                 bool                       `json:"unavailable,omitempty"`
	MemberCount                 int                        `json:"member_count,omitempty"`
	VoiceStates                 []*VoiceState              `json:"voice_states,omitempty"`
	Members                     []*GuildMember             `json:"members,omitempty"`
	Channels                    []*Channel                 `json:"channels,omitempty"`
	Presences                   []*PresenceUpdate          `json:"presences,omitempty"`
	MaxPresences                int                        `json:"max_presences,omitempty"`
	MaxMembers                  int                        `json:"max_members,omitempty"`
	VanityURLCode               string                     `json:"vanity_url_code,omitempty"`
	Description                 string                     `json:"description,omitempty"`
	Banner                      string                     `json:"banned,omitempty"`
	PremiumTier                 PremiumTierLevel           `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      Snowflake                  `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      int                        `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int                        `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen             `json:"welcome_screen,omitempty"`
	NSFWLevel                   GuildNSFWLevel             `json:"nsfw_level"`
	StageInstances              []*StageInstance           `json:"stage_instances,omitempty"`
	Stickers                    []*Sticker                 `json:"stickers"`
}

type UnavailableGuild struct {
	DiscordBaseObject
	Unavailable bool `json:"unavailable"`
}

type GuildPreview struct {
	DiscordBaseObject
	Name                     string         `json:"name"`
	Icon                     string         `json:"icon,omitempty"`
	Splash                   string         `json:"splash,omitempty"`
	DiscoverySplash          string         `json:"discovery_splash,omitempty"`
	Emojis                   []*Emoji       `json:"emojis"`
	Features                 []GuildFeature `json:"features"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
	Description              string         `json:"description,omitempty"`
}

type WidgetUser struct {
	DiscordBaseObject
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Status        string `json:"status"`
	AvatarURL     string `json:"avatar_url"`
}

type GuildWidgetJSON struct {
	DiscordBaseObject
	Name          string        `json:"name"`
	Channels      []*Channel    `json:"channels"`
	Members       []*WidgetUser `json:"members"`
	PresenceCount int           `json:"presence_count"`
}

type Ban struct {
	Reason string `json:"reason,omitempty"`
	User   *User  `json:"user"`
}

type GuildWidget struct {
	Enabled   bool      `json:"enabled"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

type Template struct {
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	UsageCount            int       `json:"usage_count"`
	CreatorID             Snowflake `json:"creator_id"`
	Creator               *User     `json:"creator"`
	CreatedAt             Time      `json:"created_at"`
	UpdatedAt             Time      `json:"updated_at"`
	SourceGuildID         Snowflake `json:"source_guild_id"`
	SerializedSourceGuild Snowflake `json:"serialized_source_guild"`
	IsDirty               bool      `json:"is_dirty"`
}

type WelcomeScreen struct {
	Description     string            `json:"description,omitempty"`
	WelcomeChannels []*WelcomeChannel `json:"welcome_channels"`
}

type WelcomeChannel struct {
	ChannelID   Snowflake `json:"channel_id"`
	Description string    `json:"description"`
	EmojiID     Snowflake `json:"emoji_id,omitempty"`
	EmojiName   string    `json:"emoji_name,omitempty"`
}

type MembershipScreening struct {
	Version     Time                        `json:"version"`
	FormFields  []*MembershipScreeningField `json:"form_fields"`
	Description string                      `json:"description,omitempty"`
}

type MembershipFieldType string

const (
	MembershipFieldTypeTerms MembershipFieldType = "TERMS"
)

type MembershipScreeningField struct {
	FieldType MembershipFieldType `json:"field_type"`
	Label     string              `json:"label"`
	Values    []string            `json:"values,omitempty"`
	Required  bool                `json:"required"`
}
