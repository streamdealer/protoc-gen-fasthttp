package middleware

import (
	"github.com/valyala/fasthttp"
)

type Funcs []Func

type Func func(next fasthttp.RequestHandler) fasthttp.RequestHandler

func (f Funcs) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	for _, m := range f {
		h = m(h)
	}

	return h
}
