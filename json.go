package interactions

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func writeJSON(ctx *fasthttp.RequestCtx, code int, body interface{}) error {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(body); err != nil {
		return err
	}

	return nil
}
