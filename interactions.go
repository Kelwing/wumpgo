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

type (
	HandlerFunc func(*objects.Interaction) *objects.InteractionResponse
)

// App is the primary interactions server
type App struct {
	Router           *fasthttprouter.Router
	server           *fasthttp.Server
	extraProps       map[string]interface{}
	propsLock        sync.RWMutex
	logger           *logrus.Logger
	restClient       *rest.Client
	commandHandler   HandlerFunc
	componentHandler HandlerFunc
}

// Create a new interactions server instance
func New(config *Config) (*App, error) {
	pubKey, err := parsePublicKey(config.PublicKey)
	if err != nil {
		return nil, err
	}

	router := fasthttprouter.New()
	a := &App{
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

	var restClient *rest.Client
	if config.RESTClient == nil {
		restClient = rest.New(&rest.Config{
			UserAgent:     "PostcordRest/1.0 (Linux) Postcord (https://github.com/Postcord)",
			Authorization: config.Token,
			Ratelimiter: rest.NewMemoryRatelimiter(&rest.MemoryConf{
				MaxRetries: 3,
			}),
		})
	} else {
		restClient = config.RESTClient
	}

	a.restClient = restClient

	return a, nil
}

// CommandHandler sets the function to handle slash command events
func (a *App) CommandHandler(handler HandlerFunc) {
	a.commandHandler = handler
}

// ComponentHandler sets the function to handle Component events.
func (a *App) ComponentHandler(handler HandlerFunc) {
	a.componentHandler = handler
}

func (a *App) requestHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	a.logger.WithField("addr", ctx.RemoteIP()).Debug("new request")
	resp, err := a.ProcessRequest(ctx.Request.Body())
	if err != nil {
		a.logger.WithError(err).Error("failed to process request: ", err)
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
func (a *App) ProcessRequest(data []byte) (resp *objects.InteractionResponse, err error) {
	var req objects.Interaction
	err = json.Unmarshal(data, &req)
	if err != nil {
		a.logger.WithError(err).Error("failed to decode request body")
		err = fmt.Errorf("failed to decode request body")
		return
	}

	a.logger.Info("received event of type ", req.Type)

	switch req.Type {
	case objects.InteractionRequestPing:
		resp = &objects.InteractionResponse{Type: objects.ResponsePong}
		return
	case objects.InteractionApplicationCommand:
		resp = a.commandHandler(&req)
	case objects.InteractionButton:
		resp = a.componentHandler(&req)
		if resp == nil {
			return &objects.InteractionResponse{
				Type: objects.ResponseDeferredMessageUpdate,
			}, nil
		}
	}

	if resp == nil {
		err = fmt.Errorf("nil response")
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
