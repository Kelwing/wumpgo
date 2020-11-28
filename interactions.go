package interactions

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"log"
	"sync"
)

type App struct {
	server     *fasthttp.Server
	commands   map[string]HandlerFunc
	extraProps map[string]interface{}
	propsLock  sync.RWMutex
}

func New(publicKey string) (*App, error) {
	pubKey, err := parsePublicKey(publicKey)
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

func (a *App) requestHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	resp, err := a.ProcessRequest(ctx.Request.Body())
	if err != nil {
		_ = writeJSON(ctx, fasthttp.StatusOK, objects.InteractionResponse{
			Type: objects.ResponseChannelMessage,
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
		return
	}

	for _, option := range ctx.Request.Data.Options {
		ctx.options[option.Name] = &CommandOption{Value: option.Value}
	}

	switch ctx.Request.Type {
	case objects.InteractionRequestPing:
		ctx = &CommandCtx{Response: &objects.InteractionResponse{Type: objects.ResponsePong}}
		return
	case objects.InteractionApplicationCommand:
		command, ok := a.commands[ctx.Request.Data.Name]
		if !ok {
			ctx.SetContent("Command doesn't have a handler.").Ephemeral()
			return
		}
		command(ctx)
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
	return a.server.ListenAndServe(fmt.Sprintf(":%d", port))
}
