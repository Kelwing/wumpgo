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
	ReactionFmt           = ReactionsBaseFmt + "/%s/%s"
	ReactionsFmt          = ReactionsBaseFmt + "/%s"
	BulkDeleteMessagesFmt = ChannelMessagesFmt + "/bulk-delete"

	// Commands
	ApplicationFmt              = BaseURL + "/applications"
	GlobalApplicationsFmt       = ApplicationFmt + "/%d/commands"
	GlobalApplicationsUpdateFmt = GlobalApplicationsFmt + "/%d"
	GuildApplicationsFmt        = ApplicationFmt + "/%d/guilds/%d/commands"
	GuildApplicationsUpdateFmt  = GuildApplicationsFmt + "/%d"
)
