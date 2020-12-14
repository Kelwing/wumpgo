package rest

const (
	// Discord API Base URL
	BaseURL = "https://discord.com/api/v8"

	JsonContentType = "application/json"

	// Channels
	ChannelBaseFmt        = BaseURL + "/channels/%d"
	ChannelMessagesFmt    = ChannelBaseFmt + "/messages"
	ChannelMessageFmt     = ChannelMessagesFmt + "/%d"
	CrosspostMessageFmt   = ChannelMessageFmt + "/crosspost"
	ReactionsBaseFmt      = ChannelMessageFmt + "/reactions"
	ReactionsFmt          = ReactionsBaseFmt + "/%s"
	ReactionFmt           = ReactionsFmt + "/%s"
	ReactionUserFmt       = ReactionsFmt + "/%d"
	BulkDeleteMessagesFmt = ChannelMessagesFmt + "/bulk-delete"
	ChannelPermissionsFmt = ChannelBaseFmt + "/permissions/%d"
	ChannelInvitesFmt     = ChannelBaseFmt + "/invites"
	ChannelPinsFmt        = ChannelBaseFmt + "/pins"
	ChannelPinnedFmt      = ChannelPinsFmt + "/%d"
	ChannelFollowersFmt   = ChannelBaseFmt + "/followers"
	ChannelTypingFmt      = ChannelBaseFmt + "/typing"

	// Commands
	ApplicationFmt              = BaseURL + "/applications"
	GlobalApplicationsFmt       = ApplicationFmt + "/%d/commands"
	GlobalApplicationsUpdateFmt = GlobalApplicationsFmt + "/%d"
	GuildApplicationsFmt        = ApplicationFmt + "/%d/guilds/%d/commands"
	GuildApplicationsUpdateFmt  = GuildApplicationsFmt + "/%d"

	// Guilds
	GuildBaseFmt      = BaseURL + "/guilds/%d"
	GuildAuditLogsFmt = GuildBaseFmt + "/audit-logs"
)
