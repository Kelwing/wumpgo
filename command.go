package router

import (
	"errors"
	"github.com/Postcord/objects"
)

// Command is used to define a Discord (sub-)command. DO NOT MAKE YOURSELF! USE CommandGroup.NewCommandBuilder OR CommandRouter.NewCommandBuilder!
type Command struct {
	// Description is the description for the command.
	Description string `json:"description"`

	// AllowedMentions is used to set a command level rule on allowed mentions. If this is not nil, it overrides the last configuration.
	AllowedMentions *objects.AllowedMentions `json:"allowed_mentions"`

	// DefaultPermission is used to define if the permissions of the command are default.
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

// Execute the command.
func (c *Command) execute(exceptionHandler func(error) *objects.InteractionResponse, allowedMentions *objects.AllowedMentions, interaction *objects.Interaction, data *objects.ApplicationCommandInteractionData, options []*objects.ApplicationCommandInteractionDataOption) *objects.InteractionResponse {
	// Process the options.
	mappedOptions := map[string]interface{}{}
	for _, v := range options {
		// Find the option.
		option := findOption(v.Name, c.Options)
		if option == nil {
			// Option was provided that was not in this command.
			return exceptionHandler(NonExistentOption)
		}

		// Check the option and result match types.
		if objects.ApplicationCommandOptionType(v.Type) != option.OptionType {
			return exceptionHandler(MismatchedOption)
		}

		// Check what the type is.
		switch option.OptionType {
		case objects.TypeChannel:
			mappedOptions[option.Name] = &ResolvableChannel{
				id:   v.Value.(string),
				data: data,
			}
		case objects.TypeRole:
			mappedOptions[option.Name] = &ResolvableRole{
				id:   v.Value.(string),
				data: data,
			}
		case objects.TypeUser:
			mappedOptions[option.Name] = &ResolvableUser{
				id:   v.Value.(string),
				data: data,
			}
		case objects.TypeString:
			mappedOptions[option.Name] = v.Value.(string)
		case objects.TypeInteger:
			mappedOptions[option.Name] = v.Value.(float64)
		case objects.TypeBoolean:
			mappedOptions[option.Name] = v.Value.(bool)
		case objects.TypeMentionable:
			mappedOptions[option.Name] = &ResolvableMentionable{
				id:   v.Value.(string),
				data: data,
			}
		}
	}

	// Get the allowed mentions configuration.
	if c.AllowedMentions != nil {
		allowedMentions = c.AllowedMentions
	}

	// Attempt to catch errors from here.
	defer func() {
		if errGeneric := recover(); errGeneric != nil {
			// Shouldn't try and return from defer.
			exceptionHandler(ungenericError(errGeneric))
		}
	}()

	// Create the context.
	rctx := &CommandRouterCtx{Interaction: interaction, Command: data.Name, Options: mappedOptions}

	// Run the command.
	handler := c.Function
	if handler == nil {
		// don't nil crash
		handler = func(ctx *CommandRouterCtx) error { return nil }
	}
	if err := handler(rctx); err != nil {
		return exceptionHandler(err)
	}
	return rctx.buildResponse(false, exceptionHandler, allowedMentions)
}
