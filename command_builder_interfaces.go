package router

import "github.com/Postcord/objects"

// TextCommandBuilder is used to define a builder for a Command object where the type is a text command.
type TextCommandBuilder interface {
	textCommandOptions

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
	subCommandOptions

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
	commandOptions

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
