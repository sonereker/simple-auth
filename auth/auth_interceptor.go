package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authInterceptor struct {
	authManager   *AuthManager
	publicMethods map[string]bool
}

//NewAuthInterceptor creates a new NewAuthInterceptor instance
func NewAuthInterceptor(am *AuthManager, publicMethods map[string]bool) *authInterceptor {
	return &authInterceptor{
		authManager:   am,
		publicMethods: publicMethods,
	}
}

//Unary creates a new UnaryServerInterceptor
func (interceptor *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *authInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	_, ok := interceptor.publicMethods[method]
	if ok {
		return nil, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.authManager.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	ctx = context.WithValue(ctx, "id", claims.ID)

	return ctx, nil
}
