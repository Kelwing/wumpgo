package interactions

import (
	"fmt"

	"github.com/Postcord/objects"
)

type CommandOptionBuilder interface {
	AddOption(objects.ApplicationCommandOption)
	GetPrefix() string
}

type CommandBuilder struct {
	cmd     *objects.ApplicationCommand
	handler HandlerFunc
	router  *CommandRouter
}

func NewCommand(router *CommandRouter) *CommandBuilder {
	return &CommandBuilder{
		cmd: &objects.ApplicationCommand{
			Options: make([]objects.ApplicationCommandOption, 0),
		},
		router: router,
	}
}

func (b *CommandBuilder) Name(name string) *CommandBuilder {
	b.cmd.Name = name
	return b
}

func (b *CommandBuilder) Description(desc string) *CommandBuilder {
	b.cmd.Description = desc
	return b
}

func (b *CommandBuilder) AddOption(option objects.ApplicationCommandOption) {
	b.cmd.Options = append(b.cmd.Options, option)
}

func (b *CommandBuilder) DefaultPermissions() *CommandBuilder {
	b.cmd.DefaultPermission = true
	return b
}

func (b *CommandBuilder) Handler(handler HandlerFunc) *CommandBuilder {
	b.handler = handler
	return b
}

func (b *CommandBuilder) Build() {
	b.router.AddCommand(b.cmd.Name, b.handler)
}

func (b *CommandBuilder) GetPrefix() string {
	return b.cmd.Name
}

type OptionBuilder struct {
	option *objects.ApplicationCommandOption
	cmd    CommandOptionBuilder
	router *CommandRouter
}

func NewOption(cmd CommandOptionBuilder, router *CommandRouter) *OptionBuilder {
	return &OptionBuilder{
		cmd:    cmd,
		option: &objects.ApplicationCommandOption{},
		router: router,
	}
}

func (b *OptionBuilder) Name(name string) *OptionBuilder {
	b.option.Name = name
	return b
}

func (b *OptionBuilder) Description(desc string) *OptionBuilder {
	b.option.Description = desc
	return b
}

func (b *OptionBuilder) Required() *OptionBuilder {
	b.option.Required = true
	return b
}

func (b *OptionBuilder) AddOption(option objects.ApplicationCommandOption) {
	b.option.Options = append(b.option.Options, option)
}

func (b *OptionBuilder) AddChoice(choice objects.ApplicationCommandOptionChoice) *OptionBuilder {
	b.option.Choices = append(b.option.Choices, choice)
	return b
}

func (b *OptionBuilder) Type(t objects.ApplicationCommandOptionType) *OptionBuilder {
	b.option.OptionType = t
	return b
}

func (b *OptionBuilder) GetPrefix() string {
	return fmt.Sprintf("%s/%s", b.cmd.GetPrefix(), b.option.Name)
}

func (b *OptionBuilder) SubCommand(handler HandlerFunc) *OptionBuilder {
	b.option.OptionType = objects.TypeSubCommand
	b.router.AddCommand(b.GetPrefix(), handler)
	return b
}

func (b *OptionBuilder) Build() {
	b.cmd.AddOption(*b.option)
}
