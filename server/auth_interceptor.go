package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//AuthInterceptor is the Interceptor struct for server
type AuthInterceptor struct {
	authManager   *AuthManager
	publicMethods map[string]bool
}

//NewAuthInterceptor creates a new NewAuthInterceptor instance
func NewAuthInterceptor(am *AuthManager, publicMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{
		authManager:   am,
		publicMethods: publicMethods,
	}
}

//Unary creates a new UnaryServerInterceptor
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

//UserIDKey is the struct used for passing User ID with context
type UserIDKey struct{}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
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

	ctx = context.WithValue(ctx, UserIDKey{}, claims.ID)

	return ctx, nil
}
