package objects

type (
	ApplicationCommandOptionType int
	InteractionType              int
	ResponseType                 int
)

type HandlerFunc func(data *Interaction) *InteractionResponse

const (
	TypeSubCommand ApplicationCommandOptionType = iota + 1
	TypeSubCommandGroup
	TypeString
	TypeInteger
	TypeBoolean
	TypeUser
	TypeChannel
	TypeRole
)

// Interaction types
const (
	InteractionRequestPing InteractionType = iota + 1
	InteractionApplicationCommand
)

// Response types
const (
	ResponsePong ResponseType = iota + 1
	ResponseAcknowledge
	ResponseChannelMessage
	ResponseChannelMessageWithSource
)

// Response flags
const (
	ResponseFlagNormal    = 0
	ResponseFlagEphemeral = 1 << 6
)

type ApplicationCommand struct {
	ID            Snowflake                  `json:"id,omitempty"`
	ApplicationID Snowflake                  `json:"application_id,omitempty"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Options       []ApplicationCommandOption `json:"options"`
	Handler       HandlerFunc                `json:"-"`
}

type ApplicationCommandOption struct {
	OptionType  ApplicationCommandOptionType     `json:"type"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Default     bool                             `json:"default"`
	Required    bool                             `json:"required"`
	Choices     []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options     []ApplicationCommandOption       `json:"options,omitempty"`
}

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

type ApplicationCommandInteractionData struct {
	ID      Snowflake                                 `json:"id"`
	Name    string                                    `json:"name"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
}

type Interaction struct {
	ID        Snowflake                          `json:"id"`
	Type      InteractionType                    `json:"type"`
	Data      *ApplicationCommandInteractionData `json:"data,omitempty"`
	GuildID   Snowflake                          `json:"guild_id"`
	ChannelID Snowflake                          `json:"channel_id"`
	Member    GuildMember                        `json:"member"`
	Token     string                             `json:"token"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int              `json:"flags"`
}

type InteractionResponse struct {
	Type ResponseType                               `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}
