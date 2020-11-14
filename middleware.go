package interactions

import (
	"crypto/ed25519"
	"github.com/Postcord/objects"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"net/http"
)

func verifyMiddleware(h fasthttprouter.Handle, key ed25519.PublicKey) fasthttprouter.Handle {
	return func(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		signature := ctx.Request.Header.Peek("x-signature-ed25519")
		if verifyMessage(ctx.PostBody(), string(signature), key) {
			h(ctx, ps)
			return
		}

		_ = writeJSON(ctx, http.StatusOK, objects.InteractionResponse{
			Type: objects.ResponseChannelMessage,
			Data: &objects.InteractionApplicationCommandCallbackData{
				Content: "An unknown error occurred",
				Flags:   objects.ResponseFlagEphemeral,
			},
		})
	}
}
