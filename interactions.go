package interactions

import (
	"encoding/json"
	"fmt"
	"github.com/Postcord/objects"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"log"
)

type App struct {
	server   *fasthttp.Server
	commands map[string]*objects.ApplicationCommand
}

func New(publicKey string) (*App, error) {
	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	router := fasthttprouter.New()
	a := &App{
		commands: make(map[string]*objects.ApplicationCommand),
		server: &fasthttp.Server{
			Handler: router.Handler,
			Name:    "Postcord",
		},
	}

	router.POST("/", verifyMiddleware(a.requestHandler, pubKey))

	return a, nil
}

func (a *App) AddCommand(command *objects.ApplicationCommand) {
	// TODO check if it exists with Discord, add if it doesn't
	a.commands[command.Name] = command
}

func (a *App) RemoveCommand(commandName string) {
	// TODO check if it exists with discord, remove if it does
	delete(a.commands, commandName)
}

func (a *App) requestHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	resp, err := a.ProcessRequest(ctx.Request.Body(), string(ctx.Request.Header.Peek("x-signature-ed25519")))
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

func (a *App) ProcessRequest(data []byte, signature string) (*objects.InteractionResponse, error) {
	payload := &objects.Interaction{}
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	var resp *objects.InteractionResponse
	switch payload.Type {
	case objects.InteractionRequestPing:
		resp = &objects.InteractionResponse{
			Type: objects.ResponsePong,
		}
	case objects.InteractionApplicationCommand:
		command, ok := a.commands[payload.Data.Name]
		if !ok {
			return &objects.InteractionResponse{
				Type: objects.ResponseChannelMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "Command doesn't have a handler.",
					Flags:   objects.ResponseFlagEphemeral,
				},
			}, nil
		}
		resp = command.Handler(payload)
	}

	return resp, nil
}

func (a *App) Run(port int) error {
	return a.server.ListenAndServe(fmt.Sprintf(":%d", port))
}
