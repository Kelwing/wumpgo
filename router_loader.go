package router

import (
	"fmt"
	"os"

	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

// ErrorHandler defines the error handler function used within Postcord.
type ErrorHandler = func(error) *objects.InteractionResponse

// Defines the builder.
type loaderBuilder struct {
	globalAllowedMentions *objects.AllowedMentions
	components            *ComponentRouter
	commands              *CommandRouter
	errHandler            ErrorHandler
	app                   HandlerAccepter
}

func (l *loaderBuilder) ComponentRouter(router *ComponentRouter) LoaderBuilder {
	l.components = router
	return l
}

func (l *loaderBuilder) ErrorHandler(cb ErrorHandler) LoaderBuilder {
	l.errHandler = cb
	return l
}

func (l *loaderBuilder) CommandRouter(router *CommandRouter) LoaderBuilder {
	l.commands = router
	return l
}

// CombinedRouter is an extension of both CommandRouter and ComponentRouter to combine the two.
// I'm personally not a huge fan of using this, but it might be appealing to some people who wish to treat it as one router.
type CombinedRouter struct {
	CommandRouter
	ComponentRouter
}

func (l *loaderBuilder) CombinedRouter(router *CombinedRouter) LoaderBuilder {
	if router == nil {
		l.components = nil
		l.commands = nil
	} else {
		l.components = &router.ComponentRouter
		l.commands = &router.CommandRouter
	}
	return l
}

func genericErrorHandler(err error) *objects.InteractionResponse {
	// Log the message.
	fmt.Println("error on route:", err)

	// Pass off to Postcord/interaction's generic handler by setting to nil.
	return nil
}

func (l *loaderBuilder) AllowedMentions(config *objects.AllowedMentions) LoaderBuilder {
	l.globalAllowedMentions = config
	return l
}

// HandlerAccepter is an interface for an object which accepts Postcord handler functions.
// In most cases, you probably want to pass through *interactions.App here.
type HandlerAccepter interface {
	ComponentHandler(handler interactions.HandlerFunc)
	CommandHandler(handler interactions.HandlerFunc)
	AutocompleteHandler(handler interactions.HandlerFunc)
	Rest() *rest.Client
}

// Defines the various bits passed through from the loader.
type loaderPassthrough struct {
	rest                  rest.RESTClient
	errHandler            ErrorHandler
	globalAllowedMentions *objects.AllowedMentions
	generateFrames        bool
}

func (l *loaderBuilder) Build(app HandlerAccepter) LoaderBuilder {
	l.app = app
	cb := l.errHandler
	if cb == nil {
		// Defines a generic error handler if the user hasn't made their own.
		cb = genericErrorHandler
	}

	generateFrames := os.Getenv("POSTCORD_GENERATE_FRAMES") == "1"

	if l.components != nil {
		// Build and load the components handler.
		handler := l.components.build(loaderPassthrough{app.Rest(), cb, l.globalAllowedMentions, generateFrames})
		app.ComponentHandler(handler)
	}

	if l.commands != nil {
		// Build and load the commands/autocomplete handler.
		commandHandler, autocompleteHandler := l.commands.build(loaderPassthrough{app.Rest(), cb, l.globalAllowedMentions, generateFrames})
		app.CommandHandler(commandHandler)
		app.AutocompleteHandler(autocompleteHandler)
	}

	return l
}

func (l *loaderBuilder) CurrentChain() (*ComponentRouter, *CommandRouter, ErrorHandler, rest.RESTClient, *objects.AllowedMentions) {
	var restClient rest.RESTClient
	if l.app != nil {
		restClient = l.app.Rest()
	}
	return l.components, l.commands, l.errHandler, restClient, l.globalAllowedMentions
}

// LoaderBuilder is the interface for a router loader builder.
type LoaderBuilder interface {
	// ComponentRouter is used to add a component router to the load process.
	ComponentRouter(*ComponentRouter) LoaderBuilder

	// CommandRouter is used to add a command router to the load process.
	CommandRouter(*CommandRouter) LoaderBuilder

	// CombinedRouter is used to add a combined router to the load process.
	CombinedRouter(router *CombinedRouter) LoaderBuilder

	// ErrorHandler is used to add an error handler to the load process.
	ErrorHandler(ErrorHandler) LoaderBuilder

	// AllowedMentions allows you to set a global allowed mentions configuration.
	AllowedMentions(*objects.AllowedMentions) LoaderBuilder

	// Build is used to execute the build.
	Build(app HandlerAccepter) LoaderBuilder

	// CurrentChain is used to get the current chain of items. Note that for obvious reasons, this is not chainable.
	// Used internally by Postcord for our testing mechanism.
	CurrentChain() (componentRouter *ComponentRouter, commandRouter *CommandRouter, errHandler ErrorHandler, restClient rest.RESTClient, allowedMentions *objects.AllowedMentions)
}

// RouterLoader is used to create a new router loader builder.
func RouterLoader() LoaderBuilder {
	return &loaderBuilder{}
}
