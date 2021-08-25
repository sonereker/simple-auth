package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/auth"
	"github.com/sonereker/simple-auth/internal"
	"github.com/sonereker/simple-auth/pb/v1"
	"github.com/sonereker/simple-auth/users"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8070", "gRPC server endpoint")
	tokenSecret        = flag.String("jwt-secret", "", "JWT secret")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s/n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := internal.NewDBConnection()
	if err != nil {
		return errors.Wrap(err, "Init Database")
	}

	err = db.AutoMigrate(&users.UserDBModel{})
	if err != nil {
		return errors.Wrap(err, "Run Migrations")
	}

	err = startGRPCServer(db)
	if err != nil {
		return errors.Wrap(err, "Start GRPC Server")
	}

	return nil
}

func startGRPCServer(db *gorm.DB) error {
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	authManager := auth.NewAuthManager(*tokenSecret)
	authInterceptor := auth.NewAuthInterceptor(authManager, publicMethods())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	}

	gs := grpc.NewServer(serverOptions...)

	usersService := users.NewUserService(db, authManager)
	pb.RegisterUsersServer(gs, usersService)

	log.Println("Running GRPC Server at " + *grpcServerEndpoint)
	if err := gs.Serve(lis); err != nil {
		return err
	}
	return nil
}

func publicMethods() map[string]bool {
	return map[string]bool{
		"/users.Users/Register": true,
	}
}
