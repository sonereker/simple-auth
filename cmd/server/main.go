package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/common"
	"github.com/sonereker/simple-auth/grpc/v1"
	"github.com/sonereker/simple-auth/users"
	ggrpc "google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
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

	err = runGRPCServer(db)
	if err != nil {
		return errors.Wrap(err, "Run gRPC Server")
	}

	err = runHTTPServer()
	if err != nil {
		return errors.Wrap(err, "Run HTTP Server")
	}
	return nil
}

func runGRPCServer(db *gorm.DB) error {
	nl, err := net.Listen("tcp", fmt.Sprintf(":%d", 9090))
	if err != nil {
		return err
	}

	grpcServer := ggrpc.NewServer()

	usersService := users.NewUsersService(db)
	grpc.RegisterUsersServer(grpcServer, usersService)

	fmt.Println("Starting GRPC Server")
	if err := grpcServer.Serve(nl); err != nil {
		return err
	}
	return nil
}

func runHTTPServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := ggrpc.Dial("localhost:9090", ggrpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	rmux := runtime.NewServeMux()
	client := grpc.NewUsersClient(conn)
	err = grpc.RegisterUsersHandlerClient(ctx, rmux, client)
	if err != nil {
		return err
	}

	log.Println("Starting REST Server")
	err = http.ListenAndServe("localhost:8080", rmux)
	if err != nil {
		return err
	}
	return nil
}
