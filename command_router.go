package router

import (
	"encoding/json"
	"errors"
	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
)

// CommandRouterCtx is used to define the commands context from the router.
type CommandRouterCtx struct {
	// Defines the response builder.
	responseBuilder

	// Defines the interaction which started this.
	*objects.Interaction

	// Command defines the command that was invoked.
	Command string `json:"command"`

	// Options is used to define any options that were set in the context. Note that if an option is unset, it will not be in the map.
	// Note that for User, Channel, Role, and Mentionable; a "*Resolvable<option type>" type is used. This will allow you to get the ID as a Snowflake, string, or attempt to get from resolved.
	Options map[string]interface{} `json:"options"`
}

// CommandGroup is a group of commands. DO NOT MAKE YOURSELF! USE CommandGroup.NewCommandGroup OR CommandRouter.NewCommandGroup!
type CommandGroup struct {
	level uint

	// Description is the description for the command group.
	Description string `json:"description"`

	// AllowedMentions is used to set a group level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions"`

	// Subcommands is a map of all of the subcommands. It is a interface{} since it can be *Command or *CommandGroup. DO NOT ADD TO THIS! USE THE ATTACHED FUNCTIONS!
	Subcommands map[string]interface{} `json:"subcommands"`
}

// GroupNestedTooDeep is thrown when the sub-command group would be nested too deep.
var GroupNestedTooDeep = errors.New("sub-command group would be nested too deep")

// NewCommandGroup is used to create a sub-command group.
func (c *CommandGroup) NewCommandGroup(name, description string) (*CommandGroup, error) {
	nextLevel := c.level + 1
	if nextLevel > 2 {
		return nil, GroupNestedTooDeep
	}
	// TODO: Validate name + description.
	g := &CommandGroup{
		level:       nextLevel,
		Description: description,
		Subcommands: map[string]interface{}{},
	}
	c.Subcommands[name] = g
	return g, nil
}

// MustNewCommandGroup calls NewCommandGroup but must succeed. If not, it will panic.
func (c *CommandGroup) MustNewCommandGroup(name, description string) *CommandGroup {
	x, err := c.NewCommandGroup(name, description)
	if err != nil {
		panic(err)
	}
	return x
}

// Defines the command builder structure.
type commandBuilder struct {
	name string
	map_ map[string]interface{}
	cmd  Command
}

func (c *commandBuilder) Description(description string) CommandBuilder {
	// TODO: Validate description
	c.cmd.Description = description
	return c
}

func (c *commandBuilder) Option(option *objects.ApplicationCommandOption) CommandBuilder {
	c.cmd.Options = append(c.cmd.Options, option)
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
	// TODO: Handle blank description.
	c.map_[c.name] = &c.cmd
	return &c.cmd, nil
}

// CommandBuilder is used to define a builder for a Command object.
type CommandBuilder interface {
	// Description is used to define the commands description.
	Description(string) CommandBuilder

	// Option is used to add a command option.
	Option(*objects.ApplicationCommandOption) CommandBuilder

	// DefaultPermission is used to define if the command should have default permissions.
	DefaultPermission() CommandBuilder

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions(*objects.AllowedMentions) CommandBuilder

	// Handler is used to add a command handler.
	Handler(func(*CommandRouterCtx) error) CommandBuilder

	// Build is used to build the command and insert it into the command router.
	Build() (*Command, error)
}

// NewCommandBuilder is used to create a builder for a *Command object.
func (c CommandGroup) NewCommandBuilder(name string) CommandBuilder {
	// TODO: Validate name
	return &commandBuilder{name: name, map_: c.Subcommands}
}

// CommandRouter is used to route commands.
type CommandRouter struct {
	roots CommandGroup
}

// NewCommandGroup is used to create a sub-command group. Works the same as CommandGroup.NewCommandGroup.
func (c *CommandRouter) NewCommandGroup(name, description string) (*CommandGroup, error) {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]interface{}{}
	}
	return c.roots.NewCommandGroup(name, description)
}

// NewCommandBuilder is used to create a builder for a *Command object.
func (c *CommandRouter) NewCommandBuilder(name string) CommandBuilder {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]interface{}{}
	}
	return &commandBuilder{name: name, map_: c.roots.Subcommands}
}

// MustNewCommandGroup calls NewCommandGroup but must succeed. If not, it will panic.
func (c *CommandRouter) MustNewCommandGroup(name, description string) *CommandGroup {
	x, err := c.NewCommandGroup(name, description)
	if err != nil {
		panic(err)
	}
	return x
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CommandRouter) MarshalJSON() ([]byte, error) {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]interface{}{}
	}
	return json.Marshal(c.roots.Subcommands)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CommandGroup) UnmarshalJSON(b []byte) error {
	res := struct {
		Description string                     `json:"description"`
		Subcommands map[string]json.RawMessage `json:"subcommands"`
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	generic := make(map[string]interface{})
	for k, v := range res.Subcommands {
		var cmd Command
		if err := json.Unmarshal(v, &cmd); err == nil {
			generic[k] = &cmd
		} else {
			var group CommandGroup
			if err := json.Unmarshal(v, &group); err != nil {
				return err
			}
			generic[k] = &group
		}
	}
	*c = CommandGroup{Description: res.Description, Subcommands: generic}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CommandRouter) UnmarshalJSON(b []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	generic := make(map[string]interface{})
	for k, v := range m {
		var cmd Command
		if err := json.Unmarshal(v, &cmd); err == nil {
			generic[k] = &cmd
		} else {
			var group CommandGroup
			if err := json.Unmarshal(v, &group); err != nil {
				return err
			}
			generic[k] = &group
		}
	}
	*c = CommandRouter{roots: CommandGroup{Subcommands: generic}}
	return nil
}

// CommandIsSubcommand is thrown when the router expects a command group and gets a command.
var CommandIsSubcommand = errors.New("expected *CommandGroup, found *Command")

// CommandIsSubcommand is thrown when the router expects a command and gets a command group.
var CommandIsNotSubcommand = errors.New("expected *Command, found *CommandGroup")

// CommandDoesNotExist is thrown when the command specified does not exist.
var CommandDoesNotExist = errors.New("the command does not exist")

// GroupDoesNotExist is thrown when the group specified does not exist.
var GroupDoesNotExist = errors.New("the group does not exist")

// Execute the group.
func (c *CommandGroup) execute(exceptionHandler func(error) *objects.InteractionResponse, allowedMentions *objects.AllowedMentions, interaction *objects.Interaction, data *objects.ApplicationCommandInteractionData, nextLevel *objects.ApplicationCommandInteractionDataOption) *objects.InteractionResponse {
	if len(data.Options) != 1 {
		// data.Options must be 1 here. A valid response will just contain the next node down the tree.
		return exceptionHandler(CommandIsNotSubcommand)
	}

	// Do a switch on the type.
	switch objects.ApplicationCommandOptionType(nextLevel.Type) {
	case objects.TypeSubCommand:
		// Expect a sub-command in the map and handle accordingly.
		cmdIface, ok := c.Subcommands[nextLevel.Name]
		if !ok {
			// The command does not exist.
			return exceptionHandler(CommandDoesNotExist)
		}
		cmd, ok := cmdIface.(*Command)
		if !ok {
			// Not a command.
			return exceptionHandler(CommandIsSubcommand)
		}
		if c.AllowedMentions != nil {
			allowedMentions = c.AllowedMentions
		}
		return cmd.execute(exceptionHandler, allowedMentions, interaction, data, nextLevel.Options)
	case objects.TypeSubCommandGroup:
		// Expect a group in the map and handle accordingly.
		cmdIface, ok := c.Subcommands[nextLevel.Name]
		if !ok {
			// The group does not exist.
			return exceptionHandler(GroupDoesNotExist)
		}
		group, ok := cmdIface.(*CommandGroup)
		if !ok {
			// Not a group.
			return exceptionHandler(CommandIsSubcommand)
		}
		if c.AllowedMentions != nil {
			allowedMentions = c.AllowedMentions
		}
		return group.execute(exceptionHandler, allowedMentions, interaction, data, nextLevel.Options[0])
	default:
		// This is just a random argument.
		return exceptionHandler(CommandIsNotSubcommand)
	}
}

// Used to build the component router by the parent.
func (c *CommandRouter) build(exceptionHandler func(error) *objects.InteractionResponse, globalAllowedMentions *objects.AllowedMentions) interactions.CommandHandlerFunc {
	baseAllowedMentions := globalAllowedMentions
	if c.roots.AllowedMentions != nil {
		baseAllowedMentions = c.roots.AllowedMentions
	}
	return func(interaction *objects.Interaction) *objects.InteractionResponse {
		// Parse the data JSON.
		var data objects.ApplicationCommandInteractionData
		if err := json.Unmarshal(interaction.Data, &data); err != nil {
			return exceptionHandler(err)
		}

		// Route the command.
		m := c.roots.Subcommands
		if m == nil {
			m = map[string]interface{}{}
		}
		cmd, ok := m[data.Name]
		if !ok {
			// Not a command.
			return nil
		}
		switch x := cmd.(type) {
		case *Command:
			// Just go ahead and call execute. That will handle the option checking anyway.
			return x.execute(exceptionHandler, baseAllowedMentions, interaction, &data, data.Options)
		case *CommandGroup:
			if len(data.Options) != 1 {
				// data.Options must be 1 here. A valid response will just contain the next node down the tree.
				return exceptionHandler(CommandIsNotSubcommand)
			}

			// Figure out if we now want the command handler or the sub-command handler.
			option := data.Options[0]
			switch objects.ApplicationCommandOptionType(option.Type) {
			case objects.TypeSubCommandGroup:
				groupIface, ok := x.Subcommands[option.Name]
				if !ok {
					// The group does not exist.
					return exceptionHandler(GroupDoesNotExist)
				}
				group, ok := groupIface.(*CommandGroup)
				if !ok {
					// Not a group.
					return exceptionHandler(CommandIsNotSubcommand)
				}
				return group.execute(exceptionHandler, baseAllowedMentions, interaction, &data, option)
			case objects.TypeSubCommand:
				cmdIface, ok := x.Subcommands[option.Name]
				if !ok {
					// The command does not exist.
					return exceptionHandler(CommandDoesNotExist)
				}
				cmd, ok := cmdIface.(*Command)
				if !ok {
					// Not a command.
					return exceptionHandler(CommandIsSubcommand)
				}
				return cmd.execute(exceptionHandler, baseAllowedMentions, interaction, &data, option.Options)
			default:
				// Not a command.
				return exceptionHandler(CommandDoesNotExist)
			}
		default:
			panic("internal error - unknown root command type")
		}
	}
}

// Get the options for a command or category.
func getOptions(cmdOrCat interface{}) []objects.ApplicationCommandOption {
	switch x := cmdOrCat.(type) {
	case *Command:
		unptr := make([]objects.ApplicationCommandOption, len(x.Options))
		for i, v := range x.Options {
			unptr[i] = *v
		}
		return unptr
	case *CommandGroup:
		cmds := make([]objects.ApplicationCommandOption, len(x.Subcommands))
		i := 0
		for k, v := range x.Subcommands {
			// Create a option based on the sub-command.
			processCommand := func(cmdName, description string, default_ bool, options []*objects.ApplicationCommandOption) objects.ApplicationCommandOption {
				if description == "" {
					description = "No description provided."
				}
				unptr := make([]objects.ApplicationCommandOption, len(options))
				for i, v := range options {
					unptr[i] = *v
				}
				return objects.ApplicationCommandOption{
					OptionType:  objects.TypeSubCommand,
					Name:        cmdName,
					Default:     default_,
					Description: description,
					Options:     unptr,
				}
			}
			switch y := v.(type) {
			case *Command:
				// Create a sub-command.
				cmds[i] = processCommand(k, y.Description, y.DefaultPermission, y.Options)
			case *CommandGroup:
				// Do some incredibly mind spiralling shit.
				description := y.Description
				if description == "" {
					description = "No description provided."
				}
				children := make([]objects.ApplicationCommandOption, len(y.Subcommands))
				childrenIndex := 0
				for k, v := range y.Subcommands {
					switch x := v.(type) {
					case *Command:
						children[childrenIndex] = processCommand(k, x.Description, x.DefaultPermission, x.Options)
					case *CommandGroup:
						description := x.Description
						if description == "" {
							description = "No description provided."
						}
						children[childrenIndex] = objects.ApplicationCommandOption{
							OptionType:  objects.TypeSubCommandGroup,
							Name:        k,
							Description: description,
							Options:     getOptions(v),
						}
					}
					childrenIndex++
				}
				cmds[i] = objects.ApplicationCommandOption{
					OptionType:  objects.TypeSubCommandGroup,
					Name:        k,
					Description: description,
					Options:     children,
				}
			}

			// Add 1 to index.
			i++
		}
		return cmds
	default:
		panic("internal error - unknown command type")
	}
}

// FormulateDiscordCommands is used to formulate the commands in such a way that they can be uploaded to Discord.
func (c *CommandRouter) FormulateDiscordCommands() []*objects.ApplicationCommand {
	cmds := make([]*objects.ApplicationCommand, len(c.roots.Subcommands))
	i := 0
	for k, v := range c.roots.Subcommands {
		// Create the command.
		description := ""
		defaultPermission := false
		switch x := v.(type) {
		case *Command:
			description = x.Description
			defaultPermission = x.DefaultPermission
		case *CommandGroup:
			description = x.Description
		}
		if description == "" {
			description = "No description provided."
		}
		cmds[i] = &objects.ApplicationCommand{
			Name:              k,
			Description:       description,
			Options:           getOptions(v),
			DefaultPermission: defaultPermission,
		}

		// Add to the index.
		i++
	}
	return cmds
}
