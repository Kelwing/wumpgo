package interactions

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Postcord/objects"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type App struct {
	Router        *fasthttprouter.Router
	server        *fasthttp.Server
	commands      map[string]HandlerFunc
	buttonHandler ButtonHandlerFunc
	extraProps    map[string]interface{}
	propsLock     sync.RWMutex
	logger        *logrus.Logger
}

func New(config *Config) (*App, error) {
	pubKey, err := parsePublicKey(config.PublicKey)
	if err != nil {
		return nil, err
	}

	router := fasthttprouter.New()
	a := &App{
		commands: make(map[string]HandlerFunc),
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

	return a, nil
}

func (a *App) AddCommand(command *objects.ApplicationCommand, h HandlerFunc) {
	// TODO check if it exists with Discord, add if it doesn't
	a.commands[command.Name] = h
}

func (a *App) RemoveCommand(commandName string) {
	// TODO check if it exists with discord, remove if it does
	delete(a.commands, commandName)
}

func (a *App) ButtonHandler(handler ButtonHandlerFunc) {
	a.buttonHandler = handler
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

func (a *App) ProcessRequest(data []byte) (ctx *CommandCtx, err error) {
	ctx = &CommandCtx{
		options: make(map[string]*CommandOption),
	}
	err = json.Unmarshal(data, ctx)
	if err != nil {
		a.logger.WithError(err).Error("failed to decode request body")
		return
	}

	a.logger.Info("received event of type ", ctx.Request.Type)

	switch ctx.Request.Type {
	case objects.InteractionRequestPing:
		ctx = &CommandCtx{Response: &objects.InteractionResponse{Type: objects.ResponsePong}}
		return
	case objects.InteractionApplicationCommand:
		var cmdData objects.ApplicationCommandInteractionData
		err = json.Unmarshal(ctx.Request.Data, &cmdData)
		if err != nil {
			a.logger.WithError(err).Error("failed to decode command data")
			ctx.SetContent("Data structure invalid.").Ephemeral()
			return
		}
		for _, option := range cmdData.Options {
			ctx.options[option.Name] = &CommandOption{Value: option.Value}
		}
		command, ok := a.commands[cmdData.Name]
		if !ok {
			a.logger.Error(cmdData.Name, " command doesn't have a handler")
			ctx.SetContent("Command doesn't have a handler.").Ephemeral()
			return
		}
		command(ctx, &cmdData)
	case objects.InteractionButton:
		if a.buttonHandler == nil {
			a.logger.Error("got button event, but button handler not set")
			return
		}

		var buttonData objects.ApplicationComponentInteractionData

		err = json.Unmarshal(ctx.Request.Data, &buttonData)
		if err != nil {
			a.logger.WithError(err).Error("failed to decode button data")
			return
		}

		a.buttonHandler(ctx, &buttonData)
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

func (a *App) Run(port int) error {
	a.logger.Info("listening on port ", port)
	return a.server.ListenAndServe(fmt.Sprintf(":%d", port))
}
