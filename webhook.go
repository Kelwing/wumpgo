package objects

//go:generate stringer -type=WebhookType -output webhook_string.go

var _ SnowflakeObject = (*Webhook)(nil)

// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type WebhookType uint

const (
	IncomingWebhook WebhookType = iota + 1
	ChannelFollowWebhook
)

// https://discord.com/developers/docs/resources/invite#invite-object-invite-structure
type Webhook struct {
	DiscordBaseObject

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
}
