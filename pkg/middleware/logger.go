package middleware

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func RequestLogger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		method := string(ctx.Method())
		path := string(ctx.Path())

		// Log the request details
		fmt.Printf("Received request: %s %s\n", method, path)

		// Call the next handler in the chain
		next(ctx)

		// Log the time taken to process the request
		duration := time.Since(start)
		fmt.Printf("Request processed in %v\n", duration)
	}
}
