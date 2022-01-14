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
	TypeNumber
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
	InteractionComponent
	InteractionAutoComplete
	InteractionModalSubmit
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
	ResponseCommandAutocompleteResult
	ResponseModal
)

type ApplicationCommandPermissionType int

const (
	PermissionTypeRole ApplicationCommandPermissionType = iota + 1
	PermissionTypeUser
)

type ApplicationCommand struct {
	// ID is the unique id of the command
	ID Snowflake `json:"id,omitempty"`
	// Type is	the type of command, defaults 1 if not set
	Type *ApplicationCommandType `json:"type,omitempty"`
	// Application ID is the unique id of the parent application
	ApplicationID Snowflake `json:"application_id,omitempty"`
	// GuildID guild id of the command, if not global
	GuildID *Snowflake `json:"guild_id,omitempty"`
	// Name is a 1-32 character name
	Name string `json:"name"`
	// Description is a 1-100 character description for CHAT_INPUT commands, empty string for USER and MESSAGE commands
	Description string `json:"description,omitempty"`
	// Options are the parameters for the command, max 25, only valid for CHAT_INPUT commands
	Options []ApplicationCommandOption `json:"options"`
	// DefaultPermission is whether the command is enabled by default when the app is added to a guild
	DefaultPermission bool `json:"default_permission"`
	// Version is an autoincrementing version identifier updated during substantial record changes
	Version Snowflake `json:"version,omitempty"`
}

type ApplicationCommandOption struct {
	OptionType   ApplicationCommandOptionType     `json:"type"`
	Name         string                           `json:"name"`
	Description  string                           `json:"description"`
	Required     bool                             `json:"required,omitempty"`
	Choices      []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options      []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes []ChannelType                    `json:"channel_types,omitempty"`
	MinValue     json.Number                      `json:"min_value,omitempty"`
	MaxValue     json.Number                      `json:"max_value,omitempty"`
	Autocomplete bool                             `json:"autocomplete,omitempty"`
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
	Type    ApplicationCommandOptionType               `json:"type"`
	Name    string                                     `json:"name"`
	Value   interface{}                                `json:"value,omitempty"`
	Focused bool                                       `json:"focused,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

type ApplicationCommandInteractionData struct {
	ID       Snowflake                                  `json:"id"`
	Name     string                                     `json:"name"`
	Type     ApplicationCommandType                     `json:"type"`
	Version  Snowflake                                  `json:"version"`
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
	Locale        string          `json:"locale"`
	GuildLocale   string          `json:"guild_locale"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             bool                              `json:"tts,omitempty"`
	Content         string                            `json:"content,omitempty"`
	Embeds          []*Embed                          `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions                  `json:"allowed_mentions,omitempty"`
	Flags           MessageFlag                       `json:"flags,omitempty"`
	Components      []*Component                      `json:"components"`
	Choices         []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	// Data for modal response
	CustomID string `json:"custom_id,omitempty"`
	Title    string `json:"title,omitempty"`
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
	// ComponentTypeInputText is only usable in modals
	ComponentTypeInputText
)

type ApplicationComponentInteractionData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
	Values        []string      `json:"values,omitempty"`
}

type ApplicationModalInteractionData struct {
	CustomID   string                          `json:"custom_id"`
	Components []*InteractionResponseComponent `json:"components"`
}

type InteractionResponseComponent struct {
	Type       ComponentType                   `json:"type"`
	CustomID   string                          `json:"custom_id"`
	Value      string                          `json:"value"`
	Components []*InteractionResponseComponent `json:"components,omitempty"`
}

type ButtonStyle int

const (
	ButtonStylePrimary = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

type TextStyle int

const (
	TextStyleShort = iota + 1
	TextStyleParagraph
)

type Style int
