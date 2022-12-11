package objects

//go:generate stringer -type=ApplicationCommandOptionType,InteractionType,ResponseType,ApplicationCommandType,ButtonStyle,ComponentType,TextStyle,ApplicationCommandPermissionType -output interactions_string.go

import (
	"encoding/json"

	"wumpgo.dev/wumpgo/objects/permissions"
)

type (
	ApplicationCommandOptionType     int
	InteractionType                  int
	ResponseType                     int
	ApplicationCommandType           int
	ButtonStyle                      int
	ComponentType                    int
	TextStyle                        int
	ApplicationCommandPermissionType int
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
	TypeAttachment
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

const (
	PermissionTypeRole ApplicationCommandPermissionType = iota + 1
	PermissionTypeUser
	PermissionTypeChannel
)

const (
	ButtonStylePrimary ButtonStyle = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

const (
	TextStyleShort TextStyle = iota + 1
	TextStyleParagraph
)

const (
	ComponentTypeActionRow ComponentType = iota + 1
	ComponentTypeButton
	ComponentTypeSelectMenu
	// ComponentTypeInputText is only usable in modals
	ComponentTypeInputText
	ComponentTypeUserSelect
	ComponentTypeRoleSelect
	ComponentTypeMentionableSelect
	ComponentTypeChannelSelect
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
	// Localization dictionary for name field. Values follow the same restrictions as name
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	// Description is a 1-100 character description for CHAT_INPUT commands, empty string for USER and MESSAGE commands
	Description string `json:"description,omitempty"`
	// Localization dictionary for description field. Values follow the same restrictions as description
	DescriptionLocalizations map[string]string `json:"description_localizations,omitempty"`
	// Options are the parameters for the command, max 25, only valid for CHAT_INPUT commands
	Options []ApplicationCommandOption `json:"options"`
	// Set of permissions represented as a bit set
	DefaultPermissions *permissions.PermissionBit `json:"default_member_permissions,omitempty"`
	// Indicates whether the command is available in DMs with the app, only for globally-scoped commands. By default, commands are visible.
	AllowUseInDMs *bool `json:"dm_permission,omitempty"`
	// DefaultPermission is whether the command is enabled by default when the app is added to a guild
	DefaultPermission *bool `json:"default_permission,omitempty"`
	// Version is an autoincrementing version identifier updated during substantial record changes
	Version Snowflake `json:"version,omitempty"`
}

type ApplicationCommandOption struct {
	OptionType               ApplicationCommandOptionType     `json:"type"`
	Name                     string                           `json:"name"`
	NameLocalizations        map[string]string                `json:"name_localizations,omitempty"`
	Description              string                           `json:"description"`
	DescriptionLocalizations map[string]string                `json:"description_localizations,omitempty"`
	Required                 bool                             `json:"required,omitempty"`
	Choices                  []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []ChannelType                    `json:"channel_types,omitempty"`
	MinValue                 json.Number                      `json:"min_value,omitempty"`
	MaxValue                 json.Number                      `json:"max_value,omitempty"`
	MinLength                int64                            `json:"min_length,omitempty"`
	MaxLength                int64                            `json:"max_length,omitempty"`
	Autocomplete             bool                             `json:"autocomplete,omitempty"`
}

type ApplicationCommandOptionChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"`
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

type ApplicationCommandDataOption struct {
	Type    ApplicationCommandOptionType    `json:"type"`
	Name    string                          `json:"name"`
	Value   interface{}                     `json:"value,omitempty"`
	Focused bool                            `json:"focused,omitempty"`
	Options []*ApplicationCommandDataOption `json:"options,omitempty"`
}

type ApplicationCommandData struct {
	ID       Snowflake                       `json:"id"`
	Name     string                          `json:"name"`
	Type     ApplicationCommandType          `json:"type"`
	Resolved ResolvedData                    `json:"resolved"`
	Options  []*ApplicationCommandDataOption `json:"options"`
	GuildID  Snowflake                       `json:"guild_id"`
	TargetID Snowflake                       `json:"target_id"`
}

type ResolvedData struct {
	Users       map[Snowflake]User        `json:"users"`
	Members     map[Snowflake]GuildMember `json:"members"`
	Roles       map[Snowflake]Role        `json:"roles"`
	Channels    map[Snowflake]Channel     `json:"channels"`
	Messages    map[Snowflake]Message     `json:"messages"`
	Attachments map[Snowflake]Attachment  `json:"attachments"`
}

type Interaction struct {
	ID             Snowflake                 `json:"id"`
	ApplicationID  Snowflake                 `json:"application_id"`
	Type           InteractionType           `json:"type"`
	Data           json.RawMessage           `json:"data,omitempty"`
	GuildID        Snowflake                 `json:"guild_id"`
	ChannelID      Snowflake                 `json:"channel_id"`
	Member         *GuildMember              `json:"member"`
	User           *User                     `json:"user"`
	Token          string                    `json:"token"`
	Version        int                       `json:"version,omitempty"`
	Message        *Message                  `json:"message,omitempty"`
	AppPermissions permissions.PermissionBit `json:"app_permissions"`
	Locale         string                    `json:"locale"`
	GuildLocale    string                    `json:"guild_locale"`
}

// Deprecated: InteractionApplicationCommandCallbackData is deprecated.
// Please see InteractionMessagesCallbackData, InteractionAutocompleteCallbackData,
// and InteractionModalCallbackData
type InteractionApplicationCommandCallbackData struct {
	TTS             bool                              `json:"tts,omitempty"`
	Content         string                            `json:"content,omitempty"`
	Embeds          []*Embed                          `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions                  `json:"allowed_mentions,omitempty"`
	Flags           MessageFlag                       `json:"flags,omitempty"`
	Components      []*Component                      `json:"components"`
	Attachments     []*Attachment                     `json:"attachments,omitempty"`
	Files           []*DiscordFile                    `json:"-"`
	Choices         []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	// Data for modal response
	CustomID string `json:"custom_id,omitempty"`
	Title    string `json:"title,omitempty"`
}

type InteractionMessagesCallbackData struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           MessageFlag      `json:"flags,omitempty"`
	Components      []*Component     `json:"components"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
	Files           []*DiscordFile   `json:"-"`
}

type InteractionAutocompleteCallbackData struct {
	Choices []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
}

type InteractionModalCallbackData struct {
	CustomID   string       `json:"custom_id,omitempty"`
	Title      string       `json:"title,omitempty"`
	Components []*Component `json:"components"`
}

type InteractionResponse struct {
	Type ResponseType `json:"type"`
	Data interface{}  `json:"data,omitempty"`
}

type MessageComponentData struct {
	CustomID      string        `json:"custom_id"`
	ComponentType ComponentType `json:"component_type"`
	Values        []string      `json:"values,omitempty"`
}

type ModalSubmitData struct {
	CustomID   string                          `json:"custom_id"`
	Components []*InteractionResponseComponent `json:"components"`
}

type InteractionResponseComponent struct {
	Type       ComponentType                   `json:"type"`
	CustomID   string                          `json:"custom_id"`
	Value      string                          `json:"value"`
	Components []*InteractionResponseComponent `json:"components,omitempty"`
}

type Component struct {
	Type        ComponentType    `json:"type"`
	CustomID    string           `json:"custom_id,omitempty"`
	Disabled    bool             `json:"disabled,omitempty"`
	Label       string           `json:"label,omitempty"`
	Style       ButtonStyle      `json:"style,omitempty"`
	Emoji       *Emoji           `json:"emoji,omitempty"`
	URL         string           `json:"url,omitempty"`
	Options     []*SelectOptions `json:"options,omitempty"`
	Placeholder string           `json:"placeholder,omitempty"`
	// Must be a pointer, discord assumes omitted value = 1
	MinValues  *int         `json:"min_values,omitempty"`
	MaxValues  *int         `json:"max_values,omitempty"`
	MinLength  *int         `json:"min_length,omitempty"`
	MaxLength  *int         `json:"max_length,omitempty"`
	Value      string       `json:"value,omitempty"`
	Required   bool         `json:"required,omitempty"`
	Components []*Component `json:"components,omitempty"`
}

type SelectOptions struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
	Default     bool   `json:"default"`
}
