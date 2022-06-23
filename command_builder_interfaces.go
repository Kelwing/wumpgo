package router

import "github.com/Postcord/objects"

// commandOptions is a struct that contains the options for a command.
type commandOptions[T any] interface {
	// StringOption is used to define an option of the type string. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 3 (STRING): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	StringOption(name, description string, required bool, choiceBuilder StringChoiceBuilder) T

	// IntOption is used to define an option of the type int. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 4 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	IntOption(name, description string, required bool, choiceBuilder IntChoiceBuilder) T

	// BoolOption is used to define an option of the type bool.
	// Maps to option type 5 (BOOLEAN): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	BoolOption(name, description string, required bool) T

	// UserOption is used to define an option of the type user.
	// Maps to option type 6 (USER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	UserOption(name, description string, required bool) T

	// ChannelOption is used to define an option of the type channel.
	// Maps to option type 7 (CHANNEL): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	ChannelOption(name, description string, required bool) T

	// RoleOption is used to define an option of the type role.
	// Maps to option type 8 (ROLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	RoleOption(name, description string, required bool) T

	// MentionableOption is used to define an option of the type mentionable.
	// Maps to option type 9 (MENTIONABLE): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	MentionableOption(name, description string, required bool) T

	// DoubleOption is used to define an option of the type double. Note that choices is ignored if it's nil or length 0.
	// Maps to option type 10 (INTEGER): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	DoubleOption(name, description string, required bool, choiceBuilder DoubleChoiceBuilder) T

	// AttachmentOption is used to define an option of the type attachment.
	// Maps to option type 11 (ATTACHMENT): https://discord.com/developers/docs/interactions/slash-commands#application-command-object-application-command-option-type
	AttachmentOption(name, description string, required bool) T
}

// TextCommandBuilder is used to define a builder for a Command object where the type is a text command.
type TextCommandBuilder interface {
	commandOptions[TextCommandBuilder]

	// Description is used to define the commands description.
	Description(string) TextCommandBuilder

	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() TextCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) TextCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) TextCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// SubCommandBuilder is similar to TextCommandBuilder but doesn't allow default permissions to be set.
type SubCommandBuilder interface {
	commandOptions[SubCommandBuilder]

	// Description is used to define the commands description.
	Description(string) SubCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) SubCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) SubCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// MessageCommandBuilder is used to define a builder for a Message object where the type is a user command.
type MessageCommandBuilder interface {
	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() MessageCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) MessageCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx, *objects.Message) error) MessageCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// UserCommandBuilder is used to define a builder for a Command object where the type is a user command.
type UserCommandBuilder interface {
	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() UserCommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) UserCommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx, *objects.GuildMember) error) UserCommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}

// CommandBuilder is used to define a builder for a Command object where the type isn't known.
type CommandBuilder interface {
	commandOptions[CommandBuilder]

	// Description is used to define the commands description.
	Description(string) CommandBuilder

	// TextCommand is used to define that this should be a text command builder.
	TextCommand() TextCommandBuilder

	// MessageCommand is used to define that this should be a message command builder.
	MessageCommand() MessageCommandBuilder

	// UserCommand is used to define that this should be a message command builder.
	UserCommand() UserCommandBuilder

	// DefaultPermission is used to define if the command should have default permissions. Note this does nothing if the command is in a group.
	DefaultPermission() CommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) CommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) CommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)

	// MustBuild is used to define when a command must build or panic.
	MustBuild() *Command
}
