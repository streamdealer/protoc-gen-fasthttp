package http

import (
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func RespondWith(ctx *fasthttp.RequestCtx, resp proto.Message, err error) {
	ctx.SetContentType("application/json")

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	data, _ := protojson.Marshal(resp)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(data)
}
