package router

import (
	"github.com/Postcord/objects"
	"github.com/Postcord/objects/permissions"
)

type commandBuilder struct {
	map_ map[string]interface{}
	cmd  Command
}

func (c *commandBuilder) Description(description string) CommandBuilder {
	c.cmd.Description = description
	return c
}

func (c *commandBuilder) DefaultPermissions(perms permissions.PermissionBit) CommandBuilder {
	c.cmd.DefaultPermissions = &perms
	return c
}

func (c *commandBuilder) GuildCommand() CommandBuilder {
	tmp := false
	c.cmd.UseInDMs = &tmp
	return c
}

func (c *commandBuilder) DefaultPermission() CommandBuilder {
	c.cmd.DefaultPermission = true
	return c
}

func (c *commandBuilder) AllowedMentions(config *objects.AllowedMentions) CommandBuilder {
	c.cmd.AllowedMentions = config
	return c
}

func (c *commandBuilder) Handler(handler func(*CommandRouterCtx) error) CommandBuilder {
	c.cmd.Function = handler
	return c
}

func (c *commandBuilder) Build() (*Command, error) {
	c.map_[c.cmd.Name] = &c.cmd
	return &c.cmd, nil
}

func (c *commandBuilder) MustBuild() *Command {
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
	x := &commandBuilder{map_: c.Subcommands, cmd: Command{Name: name, commandType: int(objects.CommandTypeChatInput), parent: c}}
	return subcommandBuilder{x}
}
