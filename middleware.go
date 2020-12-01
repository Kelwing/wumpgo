package interactions

import (
	"crypto/ed25519"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"log"
)

func verifyMiddleware(h fasthttprouter.Handle, key ed25519.PublicKey) fasthttprouter.Handle {
	return func(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		signature := ctx.Request.Header.Peek("x-signature-ed25519")
		log.Println("Headers:", ctx.Request.Header.String())
		log.Println("Body:", string(ctx.Request.Body()))
		if verifyMessage(ctx.PostBody(), string(signature), key) {
			log.Println("signature verification succeeded")
			h(ctx, ps)
			return
		}

		log.Println("signature verification failed")

		// Need to return an error code here since Discord sends random security checks
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
	}
}
