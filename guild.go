package objects

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
	GuildFeatureInviteSplash GuildFeature = "INVITE_SPLASH"
	GuildFeatureVIPRegions   GuildFeature = "VIP_REGIONS"
	GuildFeatureVanityURl    GuildFeature = "VANITY_URL"

	GuildFeatureVerified             GuildFeature = "VERIFIED"
	GuildFeaturePartnered            GuildFeature = "PARTNERED"
	GuildFeatureCommunity            GuildFeature = "COMMUNITY"
	GuildFeatureCommerce             GuildFeature = "COMMERCE"
	GuildFeatureNews                 GuildFeature = "NEWS"
	GuildFeatureDiscoverable         GuildFeature = "DISCOVERABLE"
	GuildFeatureFeaturable           GuildFeature = "FEATURABLE"
	GuildFeatureAnimatedIcon         GuildFeature = "ANIMATED_ICON"
	GuildFeatureBanner               GuildFeature = "BANNER"
	GuildFeatureWelcomeScreenEnabled GuildFeature = "WELCOME_SCREEN_ENABLED"
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

type Guild struct {
	ID                          Snowflake                  `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        string                     `json:"icon,omitempty"`
	IconHash                    string                     `json:"icon_hash,omitempty"`
	Splash                      string                     `json:"splash,omitempty"`
	DiscoverySplash             string                     `json:"discovery_splash,omitempty"`
	Owner                       bool                       `json:"owner,omitempty"`
	OwnerID                     Snowflake                  `json:"owner_id"`
	Permissions                 PermissionBit              `json:"permissions,omitempty"`
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
}

type GuildPreview struct {
	ID                       Snowflake      `json:"id"`
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
	ID            Snowflake `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Status        string    `json:"status"`
	AvatarURL     string    `json:"avatar_url"`
}

type GuildWidget struct {
	ID            Snowflake     `json:"id"`
	Name          string        `json:"name"`
	Channels      []*Channel    `json:"channels"`
	Members       []*WidgetUser `json:"members"`
	PresenceCount int           `json:"presence_count"`
}

type Ban struct {
	Reason string `json:"reason,omitempty"`
	User   *User  `json:"user"`
}
