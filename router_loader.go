package router

import (
	"fmt"

	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

// Defines the builder.
type loaderBuilder struct {
	globalAllowedMentions *objects.AllowedMentions
	components            *ComponentRouter
	commands              *CommandRouter
	errHandler            func(error) *objects.InteractionResponse
}

func (l *loaderBuilder) ComponentRouter(router *ComponentRouter) LoaderBuilder {
	l.components = router
	return l
}

func (l *loaderBuilder) ErrorHandler(cb func(error) *objects.InteractionResponse) LoaderBuilder {
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
	l.components = &router.ComponentRouter
	l.commands = &router.CommandRouter
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
	Rest() *rest.Client
}

func (l *loaderBuilder) Build(app HandlerAccepter) {
	cb := l.errHandler
	if cb == nil {
		// Defines a generic error handler if the user hasn't made their own.
		cb = genericErrorHandler
	}

	if l.components != nil {
		// Build and load the components handler.
		handler := l.components.build(app.Rest(), cb, l.globalAllowedMentions)
		app.ComponentHandler(handler)
	}

	if l.commands != nil {
		// Build and load the commands handler.
		app.CommandHandler(l.commands.build(app.Rest(), cb, l.globalAllowedMentions))
	}
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
	ErrorHandler(func(error) *objects.InteractionResponse) LoaderBuilder

	// AllowedMentions allows you to set a global allowed mentions configuration.
	AllowedMentions(*objects.AllowedMentions) LoaderBuilder

	// Build is used to execute the build.
	Build(app HandlerAccepter)
}

// RouterLoader is used to create a new router loader builder.
func RouterLoader() LoaderBuilder {
	return &loaderBuilder{}
}
