package middleware

import (
	"github.com/valyala/fasthttp"
)

type Middlewares []Middleware

type Middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler

// Chain applies middlewares in "adding" order (FIFO)
func (m Middlewares) Chain(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// Apply applies middlewares in LIFO order
func (f Middlewares) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	for _, m := range f {
		h = m(h)
	}

	return h
}
