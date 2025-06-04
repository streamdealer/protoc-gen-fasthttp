package middleware

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func Auth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := ctx.Request.Header.Peek("Authorization")

		// Log the request details
		fmt.Printf("Received Auth: %s\n", t)

		// Call the next handler in the chain
		next(ctx)
	}
}
