package users

import (
	"context"
	"fmt"
	"github.com/sonereker/simple-auth/client"
	"github.com/sonereker/simple-auth/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"testing"
	"time"
)

const (
	seedUserEmail = "seed-user@email.com"
	dummyPassword = "hello123"
)

var grpcServerAddr string

func init() {
	grpcServerAddr = os.Getenv("GRPC_SERVER_ADDR")

	// Register seed user
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, false, nil)
	defer conn.Close()

	request := getDummyRegistrationRequest(true)
	_, err := uc.Register(context.Background(), request)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRegisterWOAuthorizationCode(t *testing.T) {
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, false, nil)
	defer conn.Close()

	request := getDummyRegistrationRequest(false)
	response, err := uc.Register(context.Background(), request)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				t.Error("User should be able to register for a new account without an authorization code", err)
			} else {
				t.Error(err)
				t.FailNow()
			}
		}
	}

	if response.User.Email != seedUserEmail {
		t.Error("User should be able to register for a new account without an authorization code", err)
	}
}

func TestLoginWOAuthorizationCode(t *testing.T) {
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, false, nil)
	defer conn.Close()

	request := &pb.LoginRequest{Email: seedUserEmail, Password: dummyPassword}
	response, err := uc.Login(context.Background(), request)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				t.Error("User should be able to login to their account without an authorization code", err)
			} else {
				t.Error(err)
				t.FailNow()
			}
		}
	}

	if response.User.Email != seedUserEmail {
		t.Error("User should be able to register for a new account without an authorization code", err)
	}
}

func TestRegisterResponse(t *testing.T) {
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, false, nil)
	defer conn.Close()

	request := getDummyRegistrationRequest(false)
	response, err := uc.Register(context.Background(), request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if response.User.Email != seedUserEmail {
		t.Error("Registration response should contain user object", err)
	}

	if response.Token == "" {
		t.Error("Registration response should contain a token", err)
	}
}

func TestLoginResponse(t *testing.T) {
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, false, nil)
	defer conn.Close()

	request := &pb.LoginRequest{Email: seedUserEmail, Password: dummyPassword}
	response, err := uc.Login(context.Background(), request)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if response.User.Email != seedUserEmail {
		t.Error("Login response should contain user object", err)
	}

	if response.Token == "" {
		t.Error("Login response should contain a token", err)
	}
}

func TestGetCurrentUser(t *testing.T) {
	conn, uc := client.NewUserClientWithConnection(grpcServerAddr, true, &pb.LoginRequest{Email: seedUserEmail, Password: dummyPassword})
	defer conn.Close()

	current, err := uc.GetCurrent(context.Background(), &pb.EmptyParams{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if current.Email != seedUserEmail {
		t.Error("Endpoint does not seem to be returning correct user info", err)
	}
}

func getDummyRegistrationRequest(existingUser bool) *pb.RegistrationRequest {
	rnd := time.Now().UnixNano()
	var email string
	if existingUser {
		email = fmt.Sprintf("dummy-%d@email.com", rnd)
	} else {
		email = seedUserEmail
	}

	return &pb.RegistrationRequest{
		Email:    email,
		Password: dummyPassword,
	}
}
