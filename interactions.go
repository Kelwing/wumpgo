package objects

import "encoding/json"

type (
	ApplicationCommandOptionType int
	InteractionType              int
	ResponseType                 int
	ApplicationCommandType       int
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
	TypeMentionable
)

// ApplicationCommand types
const (
	CommandTypeChatInput ApplicationCommandType = iota + 1
	CommandTypeUser
	CommandTypeMessage
)

// Interaction types
const (
	InteractionRequestPing InteractionType = iota + 1
	InteractionApplicationCommand
	InteractionButton
)

// Response types
const (
	ResponsePong ResponseType = iota + 1
	_
	_
	ResponseChannelMessageWithSource
	ResponseDeferredChannelMessageWithSource
	ResponseDeferredMessageUpdate // buttons only
	ResponseUpdateMessage
)

// Response flags
const (
	ResponseFlagNormal    = 0
	ResponseFlagEphemeral = 1 << 6
)

type ApplicationCommandPermissionType int

const (
	PermissionTypeRole ApplicationCommandPermissionType = iota + 1
	PermissionTypeUser
)

type ApplicationCommand struct {
	ID                Snowflake                  `json:"id,omitempty"`
	ApplicationID     Snowflake                  `json:"application_id,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description,omitempty"`
	Options           []ApplicationCommandOption `json:"options"`
	DefaultPermission bool                       `json:"default_permission"`
	Type              *int                       `json:"type,omitempty"`
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

type GuildApplicationCommandPermissions struct {
	ID            Snowflake                       `json:"id"`
	ApplicationID Snowflake                       `json:"application_id"`
	GuildID       Snowflake                       `json:"guild_id"`
	Permissions   []ApplicationCommandPermissions `json:"permissions"`
}

type ApplicationCommandPermissions struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}

type ApplicationCommandInteractionDataOption struct {
	Type    int                                        `json:"type"`
	Name    string                                     `json:"name"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

type ApplicationCommandInteractionData struct {
	ID       Snowflake                                  `json:"id"`
	Name     string                                     `json:"name"`
	Options  []*ApplicationCommandInteractionDataOption `json:"options"`
	Resolved ApplicationCommandInteractionDataResolved  `json:"resolved"`
	TargetID Snowflake                                  `json:"target_id"`
}

type ApplicationCommandInteractionDataResolved struct {
	Users    map[Snowflake]User        `json:"users"`
	Members  map[Snowflake]GuildMember `json:"members"`
	Roles    map[Snowflake]Role        `json:"roles"`
	Channels map[Snowflake]Channel     `json:"channels"`
	Messages map[Snowflake]Message     `json:"messages"`
}

type Interaction struct {
	ID            Snowflake       `json:"id"`
	ApplicationID Snowflake       `json:"application_id"`
	Type          InteractionType `json:"type"`
	Data          json.RawMessage `json:"data,omitempty"`
	GuildID       Snowflake       `json:"guild_id"`
	ChannelID     Snowflake       `json:"channel_id"`
	Member        *GuildMember    `json:"member"`
	User          *User           `json:"user"`
	Token         string          `json:"token"`
	Message       *Message        `json:"message,omitempty"`
	Version       int             `json:"version,omitempty"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int              `json:"flags"`
	Components      []*Component     `json:"components"`
}

type InteractionResponse struct {
	Type ResponseType                               `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

type ComponentType int

const (
	ComponentTypeActionRow = iota + 1
	ComponentTypeButton
	ComponentTypeSelectMenu
)

type ApplicationComponentInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
	Values        []string      `json:"values,omitempty"`
}

type ButtonStyle int

const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)
