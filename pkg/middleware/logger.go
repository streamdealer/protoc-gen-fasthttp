package middleware

import (
	"time"

	"github.com/valyala/fasthttp"
)

func RequestLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		method := string(ctx.Method())
		path := string(ctx.Path())
		ctx.Logger().Printf("Received request: %s %s\n", method, path)

		next(ctx)

		duration := time.Since(start)
		ctx.Logger().Printf("Request processed in %v\n", duration)
	}
}
