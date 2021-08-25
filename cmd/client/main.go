package main

import (
	"context"
	"fmt"
	"github.com/sonereker/simple-auth/common/grpc"
	ggrpc "google.golang.org/grpc"
)

func main() {
	serverAddress := "localhost:8080"
	conn, _ := ggrpc.Dial(serverAddress, ggrpc.WithInsecure())
	client := grpc.NewUsersClient(conn)

	fmt.Println("Registering a new user")
	newUser := grpc.UserRequest{
		Email:    "dummy@email.com",
		Password: "hello123",
	}
	response, err := client.Register(context.Background(), &newUser)
	if err != nil {
		panic("There was an error registering new user")
	}

	fmt.Println("User successfully registered")
	fmt.Println(response)
}
