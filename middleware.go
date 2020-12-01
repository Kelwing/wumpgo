package interactions

import (
	"crypto/ed25519"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func verifyMiddleware(h fasthttprouter.Handle, key ed25519.PublicKey) fasthttprouter.Handle {
	return func(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		signature := ctx.Request.Header.Peek("x-signature-ed25519")
		if verifyMessage(ctx.PostBody(), string(signature), key) {
			h(ctx, ps)
			return
		}

		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
	}
}
