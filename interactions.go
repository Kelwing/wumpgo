package interactions

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
)
import "github.com/postcord/objects"

type App struct {
	server   *fasthttp.Server
	commands map[string]*objects.ApplicationCommand
	pubKey   ed25519.PublicKey
}

func New(publicKey string) (*App, error) {
	pubKey, err := parsePublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	a := &App{
		commands: make(map[string]*objects.ApplicationCommand),
		pubKey:   pubKey,
	}
	a.server = &fasthttp.Server{
		Handler: a.requestHandler,
		Name:    "Postcord",
	}
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

func (a *App) requestHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.URI().Path()) != "/" || !ctx.Request.Header.IsPost() {
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(map[string]string{"error": "Not found"})
		_, _ = ctx.Write(data)
		ctx.Response.Header.SetContentType("application/json")
		return
	}

	resp, err := a.ProcessRequest(ctx.Request.Body(), string(ctx.Request.Header.Peek("x-signature-ed25519")))
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(objects.InteractionResponse{
			Type: objects.ResponseChannelMessage,
			Data: &objects.InteractionApplicationCommandCallbackData{
				Content: "An unknown error occurred",
				Flags:   objects.ResponseFlagEphemeral,
			},
		})
		_, _ = ctx.Write(data)
		ctx.Response.Header.SetContentType("application/json")
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(resp)
	_, _ = ctx.Write(data)
	ctx.Response.Header.SetContentType("application/json")
}

func (a *App) ProcessRequest(data []byte, signature string) ([]byte, error) {
	if !a.verifyMessage(data, signature) {
		return nil, errors.New("message does not match signature")
	}
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
			return json.Marshal(objects.InteractionResponse{
				Type: objects.ResponseChannelMessage,
				Data: &objects.InteractionApplicationCommandCallbackData{
					Content: "Command doesn't have a handler.",
					Flags:   objects.ResponseFlagEphemeral,
				},
			})
		}
		resp = command.Handler(payload)
	}

	return json.Marshal(resp)
}

func (a *App) Run(port int) error {
	return a.server.ListenAndServe(fmt.Sprintf(":%d", port))
}
