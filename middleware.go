package interactions

import (
	"crypto/ed25519"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

func verifyMiddleware(h fasthttprouter.Handle, key ed25519.PublicKey) fasthttprouter.Handle {
	return func(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		signature := ctx.Request.Header.Peek("X-Signature-Ed25519")
		body := ctx.PostBody()
		body = append(ctx.Request.Header.Peek("X-Signature-Timestamp"), body...)
		if verifyMessage(body, string(signature), key) {
			h(ctx, ps)
			return
		}

		// Discord requires a 4xx response code to security checks
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
	}
}
