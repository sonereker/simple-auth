package client

import (
	"context"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAttachToken(t *testing.T) {
	interceptor := &AuthInterceptor{
		Token: "dummy-token",
	}

	incomingCtx := metadata.NewIncomingContext(context.Background(), nil)
	outgoingCtx := interceptor.attachToken(incomingCtx)
	md, ok := metadata.FromOutgoingContext(outgoingCtx)
	if !ok || md["authorization"][0] != interceptor.Token {
		t.Errorf("attachToken() did not attached token to context")
	}
}
