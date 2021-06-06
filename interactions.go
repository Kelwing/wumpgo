package interactions

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

// App is the primary interactions server
type App struct {
	Router           *fasthttprouter.Router
	server           *fasthttp.Server
	commands         map[string]CommandData
	componentHandler ComponentHandlerFunc
	extraProps       map[string]interface{}
	propsLock        sync.RWMutex
	logger           *logrus.Logger
	restClient       *rest.Client
	cmdRouter        *CommandRouter
}

// Create a new interactions server instance
func New(config *Config) (*App, error) {
	pubKey, err := parsePublicKey(config.PublicKey)
	if err != nil {
		return nil, err
	}

	router := fasthttprouter.New()
	a := &App{
		commands: make(map[string]CommandData),
		server: &fasthttp.Server{
			Handler: router.Handler,
			Name:    "Postcord",
		},
		extraProps: make(map[string]interface{}),
		Router:     router,
	}

	if config.Logger == nil {
		a.logger = logrus.StandardLogger()
	} else {
		a.logger = config.Logger
	}

	router.POST("/", verifyMiddleware(a.requestHandler, pubKey))

	restClient := rest.New(&rest.Config{
		Ratelimiter: rest.NewMemoryRatelimiter(&rest.MemoryConf{
			UserAgent:     "PostcordRest/1.0 (Linux) Postcord (https://github.com/Postcord)",
			Authorization: config.Token,
			MaxRetries:    3,
		}),
	})

	a.restClient = restClient

	a.cmdRouter = NewCommandRouter(a)

	return a, nil
}

// AddCommand adds a handler for a slash command
func (a *App) AddCommand(command *objects.ApplicationCommand, h HandlerFunc) {
	// TODO check if it exists with Discord, add if it doesn't
	a.commands[command.Name] = CommandData{
		Command: command,
		Handler: h,
	}
}

// RemoveCommand removes a handler by command name
func (a *App) RemoveCommand(commandName string) {
	// TODO check if it exists with discord, remove if it does
	delete(a.commands, commandName)
}

// Commands returns a list of registered commands.
// Useful for patching all commands to Discord.
func (a *App) Commands() []*objects.ApplicationCommand {
	cmds := make([]*objects.ApplicationCommand, 0)

	for _, cmd := range a.commands {
		cmds = append(cmds, cmd.Command)
	}

	return cmds
}

// ComponentHandler sets the function to handle Component events.
func (a *App) ComponentHandler(handler ComponentHandlerFunc) {
	a.componentHandler = handler
}

func (a *App) requestHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	a.logger.WithField("addr", ctx.RemoteIP()).Debug("new request")
	resp, err := a.ProcessRequest(ctx.Request.Body())
	if err != nil {
		a.logger.WithError(err).Error("failed to process request")
		_ = writeJSON(ctx, fasthttp.StatusOK, objects.InteractionResponse{
			Type: objects.ResponseChannelMessageWithSource,
			Data: &objects.InteractionApplicationCommandCallbackData{
				Content: "An unknown error occurred",
				Flags:   objects.ResponseFlagEphemeral,
			},
		})
		return
	}

	err = writeJSON(ctx, fasthttp.StatusOK, resp)
	if err != nil {
		log.Println("failed to write response: ", err)
	}
}

// ProcessRequest is used internally to process a validated request.
// It is exposed to allow users to tied Postcord in with any web framework
// of their choosing.  Ensure you only pass validated requests.
func (a *App) ProcessRequest(data []byte) (ctx *Ctx, err error) {
	ctx = &Ctx{
		app: a,
	}
	err = json.Unmarshal(data, ctx)
	if err != nil {
		a.logger.WithError(err).Error("failed to decode request body")
		return
	}

	a.logger.Info("received event of type ", ctx.Request.Type)

	switch ctx.Request.Type {
	case objects.InteractionRequestPing:
		ctx = &Ctx{Response: &objects.InteractionResponse{Type: objects.ResponsePong}}
		return
	case objects.InteractionApplicationCommand:
		if err := a.cmdRouter.Execute(ctx); err != nil {
			ctx.SetContent("Something went wrong").Ephemeral()
		}
	case objects.InteractionButton:
		if a.componentHandler == nil {
			a.logger.Error("got button event, but button handler not set")
			return
		}

		btnCtx := &ComponentCtx{
			Ctx: ctx,
		}

		err = json.Unmarshal(ctx.Request.Data, &btnCtx.Data)
		if err != nil {
			a.logger.WithError(err).Error("failed to decode button data")
			return
		}

		a.componentHandler(btnCtx)
	}

	return
}

// Get retrieves a value from the global context
func (a *App) Get(key string) (interface{}, bool) {
	a.propsLock.RLock()
	defer a.propsLock.RUnlock()
	obj, ok := a.extraProps[key]
	return obj, ok
}

// Set stores a value in the global context.  This is suitable for things like database connections.
func (a *App) Set(key string, obj interface{}) {
	a.propsLock.Lock()
	defer a.propsLock.Unlock()
	a.extraProps[key] = obj
}

// Run runs the App with a built-in fasthttp web server.  It takes a port as its only argument.
func (a *App) Run(port int) error {
	a.logger.Info("listening on port ", port)
	return a.server.ListenAndServe(fmt.Sprintf(":%d", port))
}

// Rest exposes the internal Rest client so you can make calls to the Discord API
func (a *App) Rest() *rest.Client {
	return a.restClient
}

// CmdRouter returns the command router so commands can be added
func (a *App) CmdRouter() *CommandRouter {
	return a.cmdRouter
}
