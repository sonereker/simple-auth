package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//AuthInterceptor is the client interceptor struct
type AuthInterceptor struct {
	Token string
}

//GetUnaryInterceptor returns UnaryClientInterceptor for the client
func (interceptor *AuthInterceptor) GetUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.Token)
}
