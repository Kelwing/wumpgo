package objects

//go:generate stringer -type=WebhookType -trimprefix=Webhook -output webhook_string.go

// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type WebhookType uint

const (
	WebhookIncoming WebhookType = iota + 1
	WebhookChannelFollower
	WebhookApplication
)

// https://discord.com/developers/docs/resources/invite#invite-object-invite-structure
type Webhook struct {
	ID Snowflake `json:"id"`

	// the type of the webhook
	Type WebhookType `json:"type"`

	// the guild id this webhook is for
	GuildID Snowflake `json:"guild_id,omitempty"`

	// the channel id this webhook is for
	ChannelID Snowflake `json:"channel_id"`

	// the user this webhook was created by (not returned when getting a webhook with its token)
	User *User `json:"user,omitempty"`

	// the default name of the webhook
	Name string `json:"name,omitempty"`

	// the default avatar of the webhook
	Avatar string `json:"avatar,omitempty"`

	// the secure token of the webhook (returned for Incoming Webhooks)
	Token string `json:"token,omitempty"`

	// the bot/OAuth2 application that created this webhook
	ApplicationID Snowflake `json:"application_id,omitempty"`

	// the guild of the channel that this webhook is following (returned for Channel Follower Webhooks)
	SourceGuild *Guild `json:"source_guild,omitempty"`

	// the channel that this webhook is following (returned for Channel Follower Webhooks)
	SourceChannel *Channel `json:"source_channel,omitempty"`

	// the url used for executing the webhook (returned by the webhooks OAuth2 flow)
	URL string `json:"url,omitempty"`
}
