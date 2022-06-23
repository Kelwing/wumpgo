package router

import (
	"container/list"
	"context"
	"errors"
	"strconv"

	"github.com/Postcord/objects"
	"github.com/Postcord/objects/permissions"
	"github.com/Postcord/rest"
)

// Command is used to define a Discord (sub-)command. DO NOT MAKE YOURSELF! USE CommandGroup.NewCommandBuilder OR CommandRouter.NewCommandBuilder!
type Command struct {
	// Defines the command type.
	commandType int

	// Defines the parent.
	parent *CommandGroup

	// Defines any autocomplete options. Interface can be any of the ___AutoCompleteFunc's.
	autocomplete map[string]any

	// Name is the commands name.
	Name string `json:"name"`

	// Description is the description for the command.
	Description string `json:"description"`

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions"`

	// DefaultPermissions indicates which users should be allowed to use this command based on their permissions.  Set to 0 to disable by default. (default: all allowed)
	DefaultPermissions *permissions.PermissionBit `json:"default_member_permissions,omitempty"`

	// UseInDMs determines if the command should be usable in DMs (default: true)
	UseInDMs *bool `json:"dm_permission,omitempty"`

	// DefaultPermission Indicates whether the command is enabled by default when the app is added to a guild, defaults to false
	// Deprecated: Not recommended for use as field will soon be deprecated.
	DefaultPermission bool `json:"default_permission"`

	// Options defines the options which are required for a command.
	Options []*objects.ApplicationCommandOption `json:"options"`

	// Function is used to define the command being called.
	Function func(*CommandRouterCtx) error `json:"-"`
}

// Finds the option.
func findOption(name string, options []*objects.ApplicationCommandOption) *objects.ApplicationCommandOption {
	if options == nil {
		return nil
	}
	for _, v := range options {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// NonExistentOption is thrown when an option is provided in an interaction that doesn't exist in the command.
var NonExistentOption = errors.New("interaction option doesn't exist on command")

// MismatchedOption is thrown when the option types mismatch.
var MismatchedOption = errors.New("mismatched interaction option")

// Defines the options for command execution.
type commandExecutionOptions struct {
	restClient       rest.RESTClient
	exceptionHandler ErrorHandler
	modalRouter      *ModalRouter
	allowedMentions  *objects.AllowedMentions
	interaction      *objects.Interaction
	data             *objects.ApplicationCommandInteractionData
	options          []*objects.ApplicationCommandInteractionDataOption
}

// Maps out the options.
func (c *Command) mapOptions(autocomplete bool, data *objects.ApplicationCommandInteractionData, options []*objects.ApplicationCommandInteractionDataOption, exceptionHandler ErrorHandler) (*objects.InteractionResponse, map[string]any) {
	mappedOptions := map[string]any{}
	for _, v := range options {
		// Find the option.
		option := findOption(v.Name, c.Options)
		if option == nil {
			// Option was provided that was not in this command.
			return exceptionHandler(NonExistentOption), nil
		}

		// Check the option and result match types.
		if v.Type != option.OptionType {
			if !autocomplete || v.Type != objects.TypeString {
				return exceptionHandler(MismatchedOption), nil
			}
		}

		// Check what the type is.
		switch option.OptionType {
		case objects.TypeChannel:
			mappedOptions[option.Name] = (ResolvableChannel)(resolvable[objects.Channel]{
				id:   v.Value.(string),
				data: data,
			})
		case objects.TypeRole:
			mappedOptions[option.Name] = (ResolvableRole)(resolvable[objects.Role]{
				id:   v.Value.(string),
				data: data,
			})
		case objects.TypeUser:
			mappedOptions[option.Name] = (ResolvableUser)(resolvableUser{resolvable[objects.User]{
				id:   v.Value.(string),
				data: data,
			}})
		case objects.TypeString:
			mappedOptions[option.Name] = v.Value.(string)
		case objects.TypeInteger:
			float, ok := v.Value.(float64)
			if ok {
				mappedOptions[option.Name] = int(float)
			} else {
				mappedOptions[option.Name] = v.Value
			}
		case objects.TypeBoolean:
			mappedOptions[option.Name] = v.Value.(bool)
		case objects.TypeMentionable:
			mappedOptions[option.Name] = (ResolvableMentionable)(resolvableMentionable{
				resolvable: resolvable[any]{
					id:   v.Value.(string),
					data: data,
				},
			})
		case objects.TypeNumber:
			mappedOptions[option.Name] = v.Value
		}
	}
	return nil, mappedOptions
}

// Execute the command.
func (c *Command) execute(reqCtx context.Context, opts commandExecutionOptions, middlewareList *list.List) (resp *objects.InteractionResponse) {
	// Process the options.
	var mappedOptions map[string]any
	if opts.data.TargetID != 0 {
		// Add a special case for "/target". The slash is there as a keyword.
		mappedOptions = map[string]any{}
		if _, ok := opts.data.Resolved.Messages[opts.data.TargetID]; ok {
			mappedOptions["/target"] = (ResolvableMessage)(resolvable[objects.Message]{
				id:   strconv.FormatUint(uint64(opts.data.TargetID), 10),
				data: opts.data,
			})
		} else if _, ok = opts.data.Resolved.Users[opts.data.TargetID]; ok {
			mappedOptions["/target"] = (ResolvableUser)(resolvableUser{
				resolvable: resolvable[objects.User]{
					id:   strconv.FormatUint(uint64(opts.data.TargetID), 10),
					data: opts.data,
				},
			})
		}
	} else {
		// Call the function to map options.
		var response *objects.InteractionResponse
		response, mappedOptions = c.mapOptions(false, opts.data, opts.options, opts.exceptionHandler)
		if mappedOptions == nil {
			return response
		}
	}

	// Get the allowed mentions configuration.
	if c.AllowedMentions != nil {
		opts.allowedMentions = c.AllowedMentions
	}

	// Attempt to catch errors from here.
	defer func() {
		if errGeneric := recover(); errGeneric != nil {
			resp = opts.exceptionHandler(ungenericError(errGeneric))
		}
	}()

	// Create the context.
	rctx := &CommandRouterCtx{
		globalAllowedMentions: opts.allowedMentions,
		errorHandler:          opts.exceptionHandler,
		modalRouter:           opts.modalRouter,
		Interaction:           opts.interaction,
		Context:               reqCtx,
		Command:               c,
		Options:               mappedOptions,
		RESTClient:            opts.restClient,
	}

	// Run the command.
	handler := c.Function
	if handler == nil {
		// don't nil crash
		handler = func(ctx *CommandRouterCtx) error { return nil }
	}
	if middlewareList.Len() == 0 {
		// Just call the command function.
		if err := handler(rctx); err != nil {
			return opts.exceptionHandler(err)
		}
	} else {
		// Wrap the end command function in a middleware function and push it.
		var middlewareWrapper MiddlewareFunc = func(ctx MiddlewareCtx) error {
			return handler(ctx.CommandRouterCtx)
		}
		middlewareList.PushBack(middlewareWrapper)

		// Create the middleware context.
		mctx := MiddlewareCtx{CommandRouterCtx: rctx, middlewareList: middlewareList}

		// Call the middleware chain.
		if err := mctx.Next(); err != nil {
			return opts.exceptionHandler(err)
		}
	}
	return rctx.buildResponse(false, opts.exceptionHandler, opts.allowedMentions)
}

// Groups is used to get the command groups that this belongs to.
func (c *Command) Groups() []*CommandGroup {
	s := []*CommandGroup{}
	for x := c.parent; x != nil; x = x.parent {
		s = append([]*CommandGroup{x}, s...)
	}
	return s
}
