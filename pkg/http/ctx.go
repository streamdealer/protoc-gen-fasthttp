package http

import (
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/encoding/protojson"
)

type Marshaler struct{}
type Unmarshaler struct{}

func SetMarshalerCtx(ctx *fasthttp.RequestCtx, options *protojson.MarshalOptions) {
	ctx.SetUserValue(Marshaler{}, options)
}
func MarshalerCtx(ctx *fasthttp.RequestCtx) (options *protojson.MarshalOptions) {
	options, ok := ctx.UserValue(Marshaler{}).(*protojson.MarshalOptions)
	if !ok || options == nil {
		options = &protojson.MarshalOptions{}
	}

	return options
}

func SetUnmarshalerCtx(ctx *fasthttp.RequestCtx, options *protojson.UnmarshalOptions) {
	ctx.SetUserValue(Unmarshaler{}, options)
}
func UnmarshalerCtx(ctx *fasthttp.RequestCtx) (options *protojson.UnmarshalOptions) {
	options, ok := ctx.UserValue(Unmarshaler{}).(*protojson.UnmarshalOptions)
	if !ok || options == nil {
		options = &protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}
	}

	return options
}
