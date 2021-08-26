package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/common"
	"github.com/sonereker/simple-auth/pb"
	"github.com/sonereker/simple-auth/server"
	"github.com/sonereker/simple-auth/users"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

var (
	grpcServerAddr = flag.String("grpc-server-addr", ":8070", "gRPC Server Address")
	tokenSecret    = flag.String("jwt-secret", "", "2AB89F28-0DF2-4D47-93AD-97810483C515")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s/n", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := common.NewDBConnection()
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
	lis, err := net.Listen("tcp", *grpcServerAddr)
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

	log.Println("Running GRPC Server at " + *grpcServerAddr)
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
