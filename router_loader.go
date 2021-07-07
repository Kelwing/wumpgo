package router

import (
	"fmt"
	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
)

// Defines the builder.
type loaderBuilder struct {
	globalAllowedMentions    *objects.AllowedMentions
	components               *ComponentRouter
	commands                 *CommandRouter
	errHandler               func(error) *objects.InteractionResponse
	strictComponentChecksOff bool
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

func genericErrorHandler(err error) *objects.InteractionResponse {
	// Log the message.
	fmt.Println("error on route:", err)

	// Pass off to Postcord/interaction's generic handler by setting to nil.
	return nil
}

func (l *loaderBuilder) DisableStrictComponentChecks() LoaderBuilder {
	l.strictComponentChecksOff = true
	return l
}

func (l *loaderBuilder) AllowedMentions(config *objects.AllowedMentions) LoaderBuilder {
	l.globalAllowedMentions = config
	return l
}

// HandlerAccepter is an interface for an object which accepts Postcord handler functions.
// In most cases, you probably want to pass through *interactions.App here.
type HandlerAccepter interface {
	ComponentHandler(handler interactions.ComponentHandlerFunc)
	CommandHandler(handler interactions.CommandHandlerFunc)
}

func (l *loaderBuilder) Build(app HandlerAccepter) {
	cb := l.errHandler
	if cb == nil {
		// Defines a generic error handler if the user hasn't made their own.
		cb = genericErrorHandler
	}

	var checker func(string) bool
	if l.components != nil {
		// Build and load the components handler.
		var handler interactions.ComponentHandlerFunc
		checker, handler = l.components.build(cb, l.globalAllowedMentions)
		app.ComponentHandler(handler)
	}

	if l.strictComponentChecksOff {
		// Nil the checker since we shouldn't be checking.
		checker = nil
	}

	if l.commands != nil {
		// Build and load the commands handler.
		app.CommandHandler(l.commands.build(cb, checker, l.globalAllowedMentions))
	}
}

// LoaderBuilder is the interface for a router loader builder.
type LoaderBuilder interface {
	// ComponentRouter is used to add a component router to the load process.
	ComponentRouter(*ComponentRouter) LoaderBuilder

	// CommandRouter is used to add a command router to the load process.
	CommandRouter(*CommandRouter) LoaderBuilder

	// ErrorHandler is used to add an error handler to the load process.
	ErrorHandler(func(error) *objects.InteractionResponse) LoaderBuilder

	// AllowedMentions allows you to set a global allowed mentions configuration.
	AllowedMentions(*objects.AllowedMentions) LoaderBuilder

	// DisableStrictComponentChecks is used to disable strict component checking.
	// Unless you are doing something incredibly complex with your infrastructure, you don't want to call this.
	DisableStrictComponentChecks() LoaderBuilder

	// Build is used to execute the build.
	Build(app HandlerAccepter)
}

// RouterLoader is used to create a new router loader builder.
func RouterLoader() LoaderBuilder {
	return &loaderBuilder{}
}
