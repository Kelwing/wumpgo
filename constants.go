package rest

const (
	// Discord API Base URL
	BaseURL = "https://discord.com/api"

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
)
