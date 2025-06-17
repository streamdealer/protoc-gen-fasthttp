package middleware

import (
	"github.com/valyala/fasthttp"
)

func Auth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := ctx.Request.Header.Peek("Authorization")
		ctx.Logger().Printf("Received Auth: %s\n", t)
		next(ctx)
	}
}
