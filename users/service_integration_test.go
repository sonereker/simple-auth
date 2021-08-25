package users

import (
	"context"
	"fmt"
	"github.com/sonereker/simple-auth/grpc/v1"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"testing"
)

var client grpc.UsersClient

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	conn, _ := ggrpc.Dial("localhost:8070", ggrpc.WithInsecure())
	client = grpc.NewUsersClient(conn)
}

func TestRegisterWOAuthorizationCode(t *testing.T) {
	fmt.Println("Registering a new user")
	newUser := grpc.RegistrationRequest{
		Email:    "dummy@email.com",
		Password: "hello123",
	}
	_, err := client.Register(context.Background(), &newUser)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				t.Error("User should be able to register for a new account without an authorization code", err)
			}
		}
	}
}

func TestRegisterWAuthorizationCode(t *testing.T) {
	fmt.Println("Registering a new user")
	newUser := grpc.RegistrationRequest{
		Email:    "dummy@email.com",
		Password: "hello123",
	}
	_, err := client.Register(context.Background(), &newUser)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				t.Error("User should be able to register for a new account with an invalid authorization code", err)
			}
		}
	}
}
