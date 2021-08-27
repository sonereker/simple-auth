package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/internal/config"
	"github.com/sonereker/simple-auth/internal/pb/v1"
	"github.com/sonereker/simple-auth/internal/server"
	"github.com/sonereker/simple-auth/internal/users"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

var (
	grpcServerAddr string
	tokenSecret    = flag.String("jwt-secret", "", "2AB89F28-0DF2-4D47-93AD-97810483C515")
)

func main() {
	grpcServerAddr = os.Getenv("GRPC_SERVER_ADDR")

	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s/n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := config.NewDBConnection()
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
	lis, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	authManager := server.NewAuthManager(*tokenSecret)
	authInterceptor := server.NewAuthInterceptor(authManager, publicMethods())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	}

	gs := grpc.NewServer(serverOptions...)

	usersService := users.NewUserService(db, authManager)
	pb.RegisterUserServiceServer(gs, usersService)

	log.Println("Running GRPC Server at " + grpcServerAddr)
	if err := gs.Serve(lis); err != nil {
		return err
	}
	return nil
}

func publicMethods() map[string]bool {
	return map[string]bool{
		"/users.UserService/Register": true,
		"/users.UserService/Login":    true,
	}
}
