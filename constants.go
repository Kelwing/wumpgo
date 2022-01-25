package rest

const (
	// Discord API Base URL
	BaseURL = "https://discord.com/api/v9"

	XAuditLogReasonHeader = "X-Audit-Log-Reason"

	JsonContentType = "application/json"

	// Channels
	ChannelBaseFmt                   = BaseURL + "/channels/%d"
	ChannelMessagesFmt               = ChannelBaseFmt + "/messages"
	ChannelMessageFmt                = ChannelMessagesFmt + "/%d"
	CrosspostMessageFmt              = ChannelMessageFmt + "/crosspost"
	ReactionsBaseFmt                 = ChannelMessageFmt + "/reactions"
	ReactionsFmt                     = ReactionsBaseFmt + "/%s"
	ReactionFmt                      = ReactionsFmt + "/%s"
	ReactionUserFmt                  = ReactionsFmt + "/%d"
	BulkDeleteMessagesFmt            = ChannelMessagesFmt + "/bulk-delete"
	ChannelPermissionsFmt            = ChannelBaseFmt + "/permissions/%d"
	ChannelInvitesFmt                = ChannelBaseFmt + "/invites"
	ChannelPinsFmt                   = ChannelBaseFmt + "/pins"
	ChannelPinnedFmt                 = ChannelPinsFmt + "/%d"
	ChannelFollowersFmt              = ChannelBaseFmt + "/followers"
	ChannelTypingFmt                 = ChannelBaseFmt + "/typing"
	ChannelMessageThreadsFmt         = ChannelMessageFmt + "/threads"
	ChannelThreadsFmt                = ChannelBaseFmt + "/threads"
	ChannelThreadMembersFmt          = ChannelBaseFmt + "/thread-members"
	ChannelThreadMembersUserFmt      = ChannelThreadMembersFmt + "/%d"
	ChannelThreadMembersMeFmt        = ChannelThreadMembersFmt + "/@me"
	ChannelThreadsActiveFmt          = ChannelThreadsFmt + "/active"
	ChannelThreadsArchivedFmt        = ChannelThreadsFmt + "/archived"
	ChannelThreadsArchivedPublicFmt  = ChannelThreadsArchivedFmt + "/public"
	ChannelThreadsArchivedPrivateFmt = ChannelThreadsArchivedFmt + "/private"
	ChannelUsersMeThreadsArchivedFmt = ChannelBaseFmt + "/users/@me/threads/archived/private"

	// Commands
	ApplicationFmt              = BaseURL + "/applications"
	GlobalApplicationsFmt       = ApplicationFmt + "/%d/commands"
	GlobalApplicationsUpdateFmt = GlobalApplicationsFmt + "/%d"
	GuildApplicationsFmt        = ApplicationFmt + "/%d/guilds/%d/commands"
	GuildApplicationsUpdateFmt  = GuildApplicationsFmt + "/%d"

	// Guilds
	GuildBaseFmt                      = BaseURL + "/guilds/%d"
	GuildCreateFmt                    = BaseURL + "/guilds"
	GuildChannelsFmt                  = GuildBaseFmt + "/channels"
	GuildPreviewFmt                   = GuildBaseFmt + "/preview"
	GuildAuditLogsFmt                 = GuildBaseFmt + "/audit-logs"
	GuildMembersFmt                   = GuildBaseFmt + "/members"
	GuildMemberFmt                    = GuildMembersFmt + "/%d"
	GuildMemberEditCurrentUserNickFmt = GuildMembersFmt + "/@me/nick"
	GuildBansFmt                      = GuildBaseFmt + "/bans"
	GuildBanUserFmt                   = GuildBansFmt + "/%d"
	GuildPruneFmt                     = GuildBaseFmt + "/prune"
	GuildVoiceRegionsFmt              = GuildBaseFmt + "/regions"
	GuildInvitesFmt                   = GuildBaseFmt + "/invites"
	GuildWidgetFmt                    = GuildBaseFmt + "/widget"
	GuildWidgetJSONFmt                = GuildWidgetFmt + ".json"
	GuildVanityURLFmt                 = GuildBaseFmt + "/vanity-url"
	GuildWidgetImageFmt               = GuildWidgetFmt + ".png"
	GuildMembershipScreeningFmt       = GuildBaseFmt + "/member-verification"
	GuildThreadsFmt                   = GuildBaseFmt + "/threads/active"

	// Roles
	GuildMemberRoleFmt = GuildBaseFmt + "/members/%d/roles/%d"
	GuildRolesFmt      = GuildBaseFmt + "/roles"
	GuildRoleFmt       = GuildRolesFmt + "/%d"

	// Integrations
	IntegrationsBaseFmt = GuildBaseFmt + "/integrations"
	IntegrationBaseFmt  = IntegrationsBaseFmt + "/%d"
	IntegrationSync     = IntegrationBaseFmt + "/sync"

	// Invites
	InviteFmt = BaseURL + "/invites/%s"

	// Templates
	TemplateFmt       = BaseURL + "/guilds/templates/%s"
	GuildTemplateFmt  = BaseURL + "/guilds/%d/templates"
	GuildTemplatesFmt = GuildTemplateFmt + "/%s"

	// Users
	UsersFmt           = BaseURL + "/users"
	UsersMeFmt         = UsersFmt + "/@me"
	UsersMeGuilds      = UsersMeFmt + "/guilds"
	UsersMeGuild       = UsersMeGuilds + "/%d"
	UsersMeGuildMember = UsersMeGuild + "/member"
	UsersMeChannels    = UsersMeFmt + "/channels"
	UserFmt            = UsersFmt + "/%d"
	UserConnections    = UsersMeFmt + "/connections"

	// Voice
	VoiceRegions = BaseURL + "/voice/regions"

	// Webhooks
	ChannelWebhookFmt   = BaseURL + "/channels/%d/webhooks"
	GuildWebhookFmt     = BaseURL + "/guilds/%d/webhooks"
	WebhookFmt          = BaseURL + "/webhooks/%d"
	WebhookWithTokenFmt = WebhookFmt + "/%s"
	WebhookMessageFmt   = WebhookWithTokenFmt + "/messages/%d"

	// Interactions
	CreateInteractionResponseFmt   = BaseURL + "/interactions/%d/%s/callback"
	OriginalInteractionResponseFmt = WebhookWithTokenFmt + "/messages/@original"

	// Permissions
	GuildApplicationCommandsPermissionsFmt = GuildApplicationsFmt + "/permissions"
	GuildApplicationCommandPermissionsFmt  = GuildApplicationsFmt + "/%d/permissions"

	// Gateway
	GatewayFmt    = BaseURL + "/gateway"
	GatewayBotFmt = GatewayFmt + "/bot"
)
