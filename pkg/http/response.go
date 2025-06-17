package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func RespondWith(ctx *fasthttp.RequestCtx, resp proto.Message, err error) {
	ctx.SetContentType("application/json")

	if err != nil {
		ErrorsConverter(ctx, err)
		return
	}

	data, _ := protojson.Marshal(resp)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(data)
}

var codesMap = map[codes.Code]int{
	codes.OK:                 fasthttp.StatusOK,
	codes.Canceled:           499, // Client Closed Request (non-standard)
	codes.Unknown:            fasthttp.StatusInternalServerError,
	codes.InvalidArgument:    fasthttp.StatusBadRequest,
	codes.DeadlineExceeded:   fasthttp.StatusRequestTimeout,
	codes.NotFound:           fasthttp.StatusNotFound,
	codes.AlreadyExists:      fasthttp.StatusConflict,
	codes.PermissionDenied:   fasthttp.StatusForbidden,
	codes.Unauthenticated:    fasthttp.StatusUnauthorized,
	codes.ResourceExhausted:  fasthttp.StatusTooManyRequests,
	codes.FailedPrecondition: fasthttp.StatusBadRequest,
	codes.Aborted:            fasthttp.StatusConflict,
	codes.OutOfRange:         fasthttp.StatusBadRequest,
	codes.Unimplemented:      fasthttp.StatusNotImplemented,
	codes.Internal:           fasthttp.StatusInternalServerError,
	codes.Unavailable:        fasthttp.StatusServiceUnavailable,
	codes.DataLoss:           fasthttp.StatusInternalServerError,
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func ErrorsConverter(ctx *fasthttp.RequestCtx, err error) {
	s, ok := status.FromError(err)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Unknown error")
		return
	}

	statusCode, ok := codesMap[s.Code()]
	if !ok {
		statusCode = fasthttp.StatusInternalServerError
	}

	// Build JSON response
	resp := ErrorResponse{}
	resp.Code = s.Code().String()
	resp.Message = s.Message()

	jsonBody, _ := json.Marshal(resp)

	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(jsonBody)
}
