package router

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/kelwing/wumpgo/interactions"
	"github.com/kelwing/wumpgo/objects"
	"github.com/kelwing/wumpgo/objects/permissions"
	"github.com/kelwing/wumpgo/rest"
)

// CommandRouterCtx is used to define the commands context from the router.
type CommandRouterCtx struct {
	// Defines the response builder. THIS MUST ALWAYS BE THE FIRST FIELD IN THE STRUCT.
	// SEE THE RESPONSE BUILDER FOR MORE INFORMATION.
	publicResponseBuilder[*CommandRouterCtx]

	// Defines the error handler.
	errorHandler ErrorHandler

	// Defines the modal router.
	modalRouter *ModalRouter

	// Defines the global allowed mentions configuration.
	globalAllowedMentions *objects.AllowedMentions

	// Defines the void ID generator.
	voidGenerator

	// Defines the interaction which started this.
	*objects.Interaction

	// Context is a context.Context passed from the HTTP handler.
	Context context.Context

	// Command defines the command that was invoked.
	Command *Command `json:"command"`

	// Options is used to define any options that were set in the context. Note that if an option is unset from Discord, it will not be in the map.
	// Note that for User, Channel, Role, and Mentionable from Discord; a "*Resolvable<option type>" type is used. This will allow you to get the ID as a Snowflake, string, or attempt to get from resolved.
	Options map[string]any `json:"options"`

	// RESTClient is used to define the REST client.
	RESTClient rest.RESTClient `json:"rest_client"`
}

// TargetMessage is used to try and get the target message. If this was not targeted at a message, returns nil.
func (c *CommandRouterCtx) TargetMessage() *objects.Message {
	message, _ := c.Options["/target"].(ResolvableMessage)
	if message == nil {
		return nil
	}
	return message.Resolve()
}

// TargetMember is used to try and get the target member. If this was not targeted at a member, returns nil.
func (c *CommandRouterCtx) TargetMember() *objects.GuildMember {
	member, _ := c.Options["/target"].(ResolvableUser)
	if member == nil {
		return nil
	}
	return member.ResolveMember()
}

// InvalidTarget is thrown when the command target is not valid.
var InvalidTarget = errors.New("wrong or no target specified")

// Used to wrap a callback that returns a targeted message into a context only friendly format.
func messageTargetWrapper(cb func(*CommandRouterCtx, *objects.Message) error) func(*CommandRouterCtx) error {
	return func(ctx *CommandRouterCtx) error {
		message := ctx.TargetMessage()
		if message == nil {
			return InvalidTarget
		}
		return cb(ctx, message)
	}
}

// Used to wrap a callback that returns a targeted member into a context only friendly format.
func memberTargetWrapper(cb func(*CommandRouterCtx, *objects.GuildMember) error) func(*CommandRouterCtx) error {
	return func(ctx *CommandRouterCtx) error {
		member := ctx.TargetMember()
		if member == nil {
			return InvalidTarget
		}
		return cb(ctx, member)
	}
}

// CommandGroup is a group of commands. DO NOT MAKE YOURSELF! USE CommandGroup.NewCommandGroup OR CommandRouter.NewCommandGroup!
type CommandGroup struct {
	level uint

	// Defines the parent.
	parent *CommandGroup

	// Middleware defines all of the groups middleware.
	Middleware []MiddlewareFunc `json:"middleware"`

	// Description is the description for the command group.
	Description string `json:"description"`

	// DefaultPermissions indicates which users should be allowed to use this command based on their permissions.  Set to 0 to disable by default. (default: all allowed)
	DefaultPermissions *permissions.PermissionBit `json:"default_member_permissions,omitempty"`

	// UseInDMs determines if the command should be usable in DMs (default: true)
	UseInDMs *bool `json:"dm_permission,omitempty"`

	// AllowedMentions is used to set a group level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions"`

	// Subcommands is a map of all of the subcommands. It is a any since it can be *Command or *CommandGroup. DO NOT ADD TO THIS! USE THE ATTACHED FUNCTIONS!
	Subcommands map[string]any `json:"subcommands"`
}

// Use is used to add middleware to the group.
func (c *CommandGroup) Use(f MiddlewareFunc) {
	c.Middleware = append(c.Middleware, f)
}

// GroupNestedTooDeep is thrown when the sub-command group would be nested too deep.
var GroupNestedTooDeep = errors.New("sub-command group would be nested too deep")

type CommandGroupOptions struct {
	DefaultPermissions permissions.PermissionBit
	UseInDMs           bool
}

// NewCommandGroup is used to create a sub-command group.
func (c *CommandGroup) NewCommandGroup(name, description string, opts *CommandGroupOptions) (*CommandGroup, error) {
	nextLevel := c.level + 1
	if nextLevel > 2 {
		return nil, GroupNestedTooDeep
	}
	var g *CommandGroup
	if opts != nil {
		// TODO: Validate name + description.
		g = &CommandGroup{
			level:              nextLevel,
			Description:        description,
			DefaultPermissions: &opts.DefaultPermissions,
			UseInDMs:           &opts.UseInDMs,
			Subcommands:        map[string]any{},
		}
	} else {
		// TODO: Validate name + description.
		g = &CommandGroup{
			level:       nextLevel,
			Description: description,
			Subcommands: map[string]any{},
		}
	}

	g.parent = c
	c.Subcommands[name] = g
	return g, nil
}

// MustNewCommandGroup calls NewCommandGroup but must succeed. If not, it will panic.
func (c *CommandGroup) MustNewCommandGroup(name, description string, opts *CommandGroupOptions) *CommandGroup {
	x, err := c.NewCommandGroup(name, description, opts)
	if err != nil {
		panic(err)
	}
	return x
}

// CommandRouter is used to route commands.
type CommandRouter struct {
	roots      CommandGroup
	middleware []MiddlewareFunc
}

// Use is used to add middleware to the router.
func (c *CommandRouter) Use(f MiddlewareFunc) {
	c.middleware = append(c.middleware, f)
}

// NewCommandGroup is used to create a sub-command group. Works the same as CommandGroup.NewCommandGroup.
func (c *CommandRouter) NewCommandGroup(name, description string, opts *CommandGroupOptions) (*CommandGroup, error) {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]any{}
	}
	g, err := c.roots.NewCommandGroup(name, description, opts)
	if err != nil {
		return nil, err
	}
	g.parent = nil
	return g, nil
}

// MustNewCommandGroup calls NewCommandGroup but must succeed. If not, it will panic.
func (c *CommandRouter) MustNewCommandGroup(name, description string, opts *CommandGroupOptions) *CommandGroup {
	x, err := c.NewCommandGroup(name, description, opts)
	if err != nil {
		panic(err)
	}
	return x
}

// NewCommandBuilder is used to create a builder for a *Command object.
func (c *CommandRouter) NewCommandBuilder(name string) CommandBuilder {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]any{}
	}
	return &commandBuilder[CommandBuilder]{cmd: Command{Name: name}, map_: c.roots.Subcommands}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CommandRouter) MarshalJSON() ([]byte, error) {
	if c.roots.Subcommands == nil {
		c.roots.Subcommands = map[string]any{}
	}
	return json.Marshal(c.roots.Subcommands)
}

// CommandIsSubcommand is thrown when the router expects a command group and gets a command.
var CommandIsSubcommand = errors.New("expected *CommandGroup, found *Command")

// CommandIsNotSubcommand is thrown when the router expects a command and gets a command group.
var CommandIsNotSubcommand = errors.New("expected *Command, found *CommandGroup")

// CommandDoesNotExist is thrown when the command specified does not exist.
var CommandDoesNotExist = errors.New("the command does not exist")

// NoAutoCompleteFunc is thrown when Discord sends a focused argument without an autocomplete function.
var NoAutoCompleteFunc = errors.New("discord sent auto-complete for argument without auto-complete function")

// Used to define the autocomplete handler.
func (c *CommandRouter) autocompleteHandler(loader loaderPassthrough) interactions.HandlerFunc {
	return func(reqCtx context.Context, interaction *objects.Interaction) *objects.InteractionResponse {
		// Parse the data JSON.
		var rootData objects.ApplicationCommandInteractionData
		if err := json.Unmarshal(interaction.Data, &rootData); err != nil {
			return loader.errHandler(err)
		}

		// Wrap the data to let us traverse the tree easier.
		var data dataWrapper = rootDataWrapper{&rootData}
		options := rootData.Options

		// Get the map of (sub-)commands.
		m := c.roots.Subcommands
		if m == nil {
			m = map[string]any{}
		}

		// Handle the traversal.
		var cmd *Command
		route := []string{"testframes", "autocompletes"}
	cmdFor:
		for {
			// Add to the route.
			route = append(route, data.name())

			// Get the item from the map.
			cmdOrCat, ok := m[data.name()]
			if !ok {
				// No command.
				if _, ok = data.(rootDataWrapper); !ok {
					// Backwards compatibility.
					loader.errHandler(CommandDoesNotExist)
				}
				return nil
			}

			// Check the type of the item.
			switch x := cmdOrCat.(type) {
			case *Command:
				// Set the object and break.
				cmd = x
				break cmdFor
			case *CommandGroup:
				// How we handle this depends on what we are expecting.
				typeIface := data.type_()
				switch type_ := typeIface.(type) {
				case objects.ApplicationCommandOptionType:
					// Check the type of the option to make sure it is a command.
					if type_ != objects.TypeSubCommandGroup {
						// If the type is anything other than a subcommand group, that does not match with this being a group.
						return loader.errHandler(CommandIsNotSubcommand)
					}
				case objects.ApplicationCommandType:
					// If this is the case, we are in the root. Look ahead to see if we are a group or a command.
					if len(options) != 1 || (options[0].Type != objects.TypeSubCommand && options[0].Type != objects.TypeSubCommandGroup) {
						// We are not a group. We know this because a root command acting as a group can only have one
						// option which is either another group or a subcommand.
						return loader.errHandler(CommandIsNotSubcommand)
					}
				default:
					// This should never happen.
					panic("postcord internal error - unknown command Type field type")
				}

				// Set the map to the subcommands from this group.
				m = x.Subcommands

				// Get the next data and setup for the next iteration.
				nextData := options[0]
				options = nextData.Options
				data = optionDataWrapper{nextData}
			}
		}

		// Create the rest tape if this is wanted.
		r := loader.rest
		tape := tape{}
		var returnedErr string
		errHandler := loader.errHandler
		if loader.generateFrames {
			r = &restTape{
				tape: &tape,
				rest: r,
			}
			errHandler = func(err error) *objects.InteractionResponse {
				returnedErr = err.Error()
				return loader.errHandler(err)
			}
		}

		// Create the command context.
		_, mappedOptions := cmd.mapOptions(true, &rootData, options, errHandler)
		if mappedOptions == nil {
			return nil
		}
		ctx := &CommandRouterCtx{
			errorHandler: errHandler,
			Interaction:  interaction,
			Command:      cmd,
			Options:      mappedOptions,
			RESTClient:   r,
		}

		// Now we have the command, we can process the autocomplete.
		for _, v := range options {
			if v.Focused {
				// Get the autocomplete function.
				f := cmd.autocomplete[v.Name]
				if f == nil {
					loader.errHandler(NoAutoCompleteFunc)
					return nil
				}

				// Defer writing the rest tape.
				var resp *objects.InteractionResponse
				defer func() {
					// Not a resource leak. GoLand thinks it is because it is in a for loop, but it always returns.

					if loader.generateFrames {
						// Now we have all the data, we can generate the frame.
						fr := frame{interaction, tape, returnedErr, resp}
						go fr.write(route...)
					}
				}()

				// Get the options.
				var resultOptions []*objects.ApplicationCommandOptionChoice
				switch x := f.(type) {
				case StringAutoCompleteFunc:
					stringifiedOptions, err := x(ctx)
					if err != nil {
						errHandler(err)
						return nil
					}
					resultOptions = make([]*objects.ApplicationCommandOptionChoice, len(stringifiedOptions))
					for i, v := range stringifiedOptions {
						resultOptions[i] = &objects.ApplicationCommandOptionChoice{
							Name:  v.Name,
							Value: v.Value,
						}
					}
				case IntAutoCompleteFunc:
					intOptions, err := x(ctx)
					if err != nil {
						errHandler(err)
						return nil
					}
					resultOptions = make([]*objects.ApplicationCommandOptionChoice, len(intOptions))
					for i, v := range intOptions {
						resultOptions[i] = &objects.ApplicationCommandOptionChoice{
							Name:  v.Name,
							Value: v.Value,
						}
					}
				case DoubleAutoCompleteFunc:
					doubleOptions, err := x(ctx)
					if err != nil {
						errHandler(err)
						return nil
					}
					resultOptions = make([]*objects.ApplicationCommandOptionChoice, len(doubleOptions))
					for i, v := range doubleOptions {
						resultOptions[i] = &objects.ApplicationCommandOptionChoice{
							Name:  v.Name,
							Value: v.Value,
						}
					}
				default:
					panic("postcord internal error - unknown autocomplete type")
				}

				// We have successfully got the result.
				resp = &objects.InteractionResponse{
					Type: objects.ResponseCommandAutocompleteResult,
					Data: &objects.InteractionApplicationCommandCallbackData{
						Choices: resultOptions,
					},
				}
				return resp
			}
		}

		// None focused. This should never happen.
		return nil
	}
}

type dataWrapper interface {
	type_() any
	name() string
}

type rootDataWrapper struct {
	root *objects.ApplicationCommandInteractionData
}

func (r rootDataWrapper) type_() any {
	return r.root.Type
}

func (r rootDataWrapper) name() string {
	return r.root.Name
}

type optionDataWrapper struct {
	option *objects.ApplicationCommandInteractionDataOption
}

func (r optionDataWrapper) type_() any {
	return r.option.Type
}

func (r optionDataWrapper) name() string {
	return r.option.Name
}

// Used to define the command handler.
func (c *CommandRouter) commandHandler(loader loaderPassthrough) interactions.HandlerFunc {
	// Get the allowed mentions configuration.
	baseAllowedMentions := loader.globalAllowedMentions
	if c.roots.AllowedMentions != nil {
		baseAllowedMentions = c.roots.AllowedMentions
	}

	// Process the response.
	return func(reqCtx context.Context, interaction *objects.Interaction) *objects.InteractionResponse {
		// Handle middleware.
		middlewareList := list.New()
		if c.middleware != nil {
			for _, v := range c.middleware {
				middlewareList.PushBack(v)
			}
		}

		// Parse the data JSON.
		var rootData objects.ApplicationCommandInteractionData
		if err := json.Unmarshal(interaction.Data, &rootData); err != nil {
			return loader.errHandler(err)
		}

		// Defines the items changed whilst traversing the tree.
		options := rootData.Options
		allowedMentions := baseAllowedMentions
		var data dataWrapper = rootDataWrapper{&rootData}

		// Get the map of (sub-)commands.
		m := c.roots.Subcommands
		if m == nil {
			m = map[string]any{}
		}

		// Create the rest tape if this is wanted.
		r := loader.rest
		tape := tape{}
		var returnedErr string
		errHandler := loader.errHandler
		if loader.generateFrames {
			r = &restTape{
				tape: &tape,
				rest: r,
			}
			errHandler = func(err error) *objects.InteractionResponse {
				returnedErr = err.Error()
				return loader.errHandler(err)
			}
		}

		// Find the route.
		route := []string{"testframes", "commands"}
		for {
			// Add to the route.
			route = append(route, data.name())

			// Get the item from the map.
			cmdOrCat, ok := m[data.name()]
			if !ok {
				// No command.
				return nil
			}

			// Check the type of the item.
			switch x := cmdOrCat.(type) {
			case *Command:
				// In this case, we should go ahead and execute.
				resp := x.execute(reqCtx, commandExecutionOptions{
					restClient:       r,
					exceptionHandler: errHandler,
					allowedMentions:  allowedMentions,
					interaction:      interaction,
					modalRouter:      loader.modalRouter,
					data:             &rootData,
					options:          options,
				}, middlewareList)
				if loader.generateFrames {
					// Now we have all the data, we can generate the frame.
					f := frame{interaction, tape, returnedErr, resp}
					go f.write(route...)
				}
				return resp
			case *CommandGroup:
				// How we handle this depends on what we are expecting.
				typeIface := data.type_()
				switch type_ := typeIface.(type) {
				case objects.ApplicationCommandOptionType:
					// Check the type of the option to make sure it is a command.
					if type_ != objects.TypeSubCommandGroup {
						// If the type is anything other than a subcommand group, that does not match with this being a group.
						return loader.errHandler(CommandIsNotSubcommand)
					}
				case objects.ApplicationCommandType:
					// If this is the case, we are in the root. Look ahead to see if we are a group or a command.
					if len(options) != 1 || (options[0].Type != objects.TypeSubCommand && options[0].Type != objects.TypeSubCommandGroup) {
						// We are not a group. We know this because a root command acting as a group can only have one
						// option which is either another group or a subcommand.
						return loader.errHandler(CommandIsNotSubcommand)
					}
				default:
					// This should never happen.
					panic("postcord internal error - unknown command Type field type")
				}

				// Handle allowed mentions.
				if x.AllowedMentions != nil {
					allowedMentions = x.AllowedMentions
				}

				// Handle middleware.
				if x.Middleware != nil {
					for _, v := range x.Middleware {
						middlewareList.PushBack(v)
					}
				}

				// Set the map to the subcommands from this group.
				m = x.Subcommands

				// Get the next data and setup for the next iteration.
				nextData := options[0]
				options = nextData.Options
				data = optionDataWrapper{nextData}
			default:
				// This should never actually happen. We need to know right away if it does.
				panic("postcord internal error - unknown root command type")
			}
		}
	}
}

// Used to build the command router by the parent.
func (c *CommandRouter) build(loader loaderPassthrough) (interactions.HandlerFunc, interactions.HandlerFunc) {
	return c.commandHandler(loader), c.autocompleteHandler(loader)
}

// Get the options for a command or category.
func getOptions(cmdOrCat any) []objects.ApplicationCommandOption {
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
			processCommand := func(cmdName, description string, options []*objects.ApplicationCommandOption) objects.ApplicationCommandOption {
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
					Description: description,
					Options:     unptr,
				}
			}
			switch y := v.(type) {
			case *Command:
				// Create a sub-command.
				cmds[i] = processCommand(k, y.Description, y.Options)
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
						children[childrenIndex] = processCommand(k, x.Description, x.Options)
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
		panic("postcord internal error - unknown command type")
	}
}

// FormulateDiscordCommands is used to formulate the commands in such a way that they can be uploaded to Discord.
func (c *CommandRouter) FormulateDiscordCommands() []*objects.ApplicationCommand {
	cmds := make([]*objects.ApplicationCommand, len(c.roots.Subcommands))
	i := 0
	for k, v := range c.roots.Subcommands {
		// Create the command.
		description := ""
		commandType := objects.CommandTypeChatInput

		cmd := &objects.ApplicationCommand{
			Name:        k,
			Description: description,
			Options:     getOptions(v),
		}

		switch x := v.(type) {
		case *Command:
			cmd.Description = x.Description
			cmd.DefaultPermissions = x.DefaultPermissions
			cmd.AllowUseInDMs = x.UseInDMs
			if x.commandType != 0 {
				commandType = objects.ApplicationCommandType(x.commandType)
			}
		case *CommandGroup:
			cmd.Description = x.Description
			cmd.DefaultPermissions = x.DefaultPermissions
			cmd.AllowUseInDMs = x.UseInDMs
		}

		if cmd.Description == "" {
			// If the description is mandatory, set it to a none provided message.
			if commandType == objects.CommandTypeChatInput {
				cmd.Description = "No description provided."
			}
		} else if commandType != objects.CommandTypeChatInput {
			// If no description is mandatory, make sure it is unset.
			cmd.Description = ""
		}

		cmd.Type = &commandType
		cmds[i] = cmd

		// Add to the index.
		i++
	}
	return cmds
}

// Tag name for option parsing
const selectorTagName = "discord"

// Bind allows you to bind the option values to a struct for easy access
func (c *CommandRouterCtx) Bind(data any) error {
	structType := reflect.TypeOf(data)
	if structType.Kind() != reflect.Ptr {
		return errors.New("data must be a pointer")
	}

	elem := structType.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("value of data pointer must be a struct")
	}

	v := reflect.ValueOf(data).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Tag == "" {
			continue
		}

		tagValue := field.Tag.Get(selectorTagName)
		if tagValue == "" {
			continue
		}

		kind := field.Type.Kind()

		option, ok := c.Options[tagValue]
		if !ok {
			continue
		}
		optionVal := reflect.ValueOf(option)

		i := -1
		for j, v := range c.Command.Options {
			if v.Name == tagValue {
				i = j
				break
			}
		}

		if !(i < len(c.Command.Options) && c.Command.Options[i].Name == tagValue) {
			continue
		}

		fieldPointer := v.FieldByName(field.Name)
		if !fieldPointer.CanSet() {
			continue
		}

		switch c.Command.Options[i].OptionType {
		case objects.TypeString:
			if kind != reflect.String {
				return fmt.Errorf("option %s is a StringOption, but the struct type is not a string", tagValue)
			}
			fmt.Printf("setting %s to %s\n", field.Name, optionVal.String())
			fieldPointer.Set(optionVal)
		case objects.TypeInteger:
			if kind != reflect.Int && kind != reflect.Int64 {
				return fmt.Errorf("option %s is a IntegerOption, but the struct type is not an integer", tagValue)
			}
			fmt.Printf("setting %s to %d\n", field.Name, optionVal.Int())
			fieldPointer.Set(optionVal)
		case objects.TypeNumber:
			if kind != reflect.Float64 {
				return fmt.Errorf("option %s is a DoubleOption, but the struct type is not a float64", tagValue)
			}
			fmt.Printf("setting %s to %f\n", field.Name, optionVal.Float())
			fieldPointer.Set(optionVal)
		case objects.TypeBoolean:
			if kind != reflect.Bool {
				return fmt.Errorf("option %s is a BoolOption, but the struct type is not a bool", tagValue)
			}
			fmt.Printf("setting %s to %t\n", field.Name, optionVal.Bool())
			fieldPointer.Set(optionVal)
		case objects.TypeChannel, objects.TypeRole, objects.TypeUser, objects.TypeMentionable:
			if kind != reflect.Interface {
				return fmt.Errorf("option %s is a Resolvable type, but the interface type is not a Resolvable", tagValue)
			}

			if optionVal.Kind() == reflect.Ptr && optionVal.IsNil() {
				return fmt.Errorf("option %s is a Resolvable type, but the value is nil", tagValue)
			}

			if optionVal.IsZero() {
				return fmt.Errorf("option %s is a Resolvable type, but the type is not a Resolvable", tagValue)
			}

			fmt.Printf("setting %s to %s\n", field.Name, optionVal.String())
			fieldPointer.Set(optionVal)
		default:
			return fmt.Errorf("option %s has an incompatible type", tagValue)
		}
	}

	return nil
}
