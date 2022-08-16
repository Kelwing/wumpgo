package router

import (
	"github.com/Postcord/objects"
	"github.com/Postcord/objects/permissions"
)

type commandBuilder[T any] struct {
	map_ map[string]any
	cmd  Command
}

func (c *commandBuilder[T]) Description(description string) T {
	c.cmd.Description = description
	return builderWrapify(c)
}

func (c *commandBuilder[T]) DefaultPermissions(perms permissions.PermissionBit) T {
	c.cmd.DefaultPermissions = &perms
	return builderWrapify(c)
}

func (c *commandBuilder[T]) GuildCommand() T {
	tmp := false
	c.cmd.UseInDMs = &tmp
	return builderWrapify(c)
}

func (c *commandBuilder[T]) AllowedMentions(config *objects.AllowedMentions) T {
	c.cmd.AllowedMentions = config
	return builderWrapify(c)
}

func (c *commandBuilder[T]) Handler(handler func(*CommandRouterCtx) error) T {
	c.cmd.Function = handler
	return builderWrapify(c)
}

func (c *commandBuilder[T]) Build() (*Command, error) {
	c.map_[c.cmd.Name] = &c.cmd
	return &c.cmd, nil
}

func (c *commandBuilder[T]) MustBuild() *Command {
	cmd, err := c.Build()
	if err != nil {
		panic(err)
	}
	return cmd
}

func (c textCommandBuilder) Description(description string) TextCommandBuilder {
	c.commandBuilder.Description(description)
	return c
}

func (c textCommandBuilder) Handler(handler func(*CommandRouterCtx) error) TextCommandBuilder {
	c.commandBuilder.Handler(handler)
	return c
}

func (c subcommandBuilder) Description(description string) SubCommandBuilder {
	c.commandBuilder.Description(description)
	return c
}

func (c subcommandBuilder) Handler(handler func(*CommandRouterCtx) error) SubCommandBuilder {
	c.commandBuilder.Handler(handler)
	return c
}

func (c messageCommandBuilder) Handler(handler func(*CommandRouterCtx, *objects.Message) error) MessageCommandBuilder {
	c.commandBuilder.Handler(messageTargetWrapper(handler))
	return c
}

func (c userCommandBuilder) Handler(handler func(*CommandRouterCtx, *objects.GuildMember) error) UserCommandBuilder {
	c.commandBuilder.Handler(memberTargetWrapper(handler))
	return c
}

// NewCommandBuilder is used to create a builder for a *Command object.
func (c *CommandGroup) NewCommandBuilder(name string) SubCommandBuilder {
	x := &commandBuilder[SubCommandBuilder]{map_: c.Subcommands, cmd: Command{Name: name, commandType: int(objects.CommandTypeChatInput), parent: c}}
	return subcommandBuilder{x}
}
