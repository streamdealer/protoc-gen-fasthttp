package interceptor

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Chain []UnaryServerInterceptor

type UnaryServerInfo struct {
	// Server is the service implementation the user provides. This is read-only.
	Server proto.Message
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
}

type UnaryHandler func(context.Context, proto.Message) (proto.Message, error)

type UnaryServerInterceptor func(
	ctx context.Context,
	req proto.Message,
	info *UnaryServerInfo,
	handler UnaryHandler,
) (
	resp proto.Message,
	err error,
)

func (interceptors Chain) Apply() UnaryServerInterceptor {
	n := len(interceptors)

	if n == 0 {
		return func(
			ctx context.Context,
			req proto.Message,
			info *UnaryServerInfo,
			handler UnaryHandler,
		) (proto.Message, error) {
			return handler(ctx, req)
		}
	}

	return func(
		ctx context.Context,
		req proto.Message,
		info *UnaryServerInfo,
		handler UnaryHandler,
	) (proto.Message, error) {
		currHandler := handler
		for i := n - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := currHandler
			currHandler = func(currentCtx context.Context, currentReq proto.Message) (proto.Message, error) {
				return interceptor(currentCtx, currentReq, info, next)
			}
		}
		return currHandler(ctx, req)
	}
}
