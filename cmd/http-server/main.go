package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/common"
	"github.com/sonereker/simple-auth/grpc/v1"
	"github.com/sonereker/simple-auth/users"
	ggrpc "google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8070", "gRPC server endpoint")
	httpServerEndpoint = flag.String("http-server-endpoint", "localhost:8080", "HTTP server endpoint")
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

	err = startHTTPServer()
	if err != nil {
		return errors.Wrap(err, "Start HTTP Server")
	}
	return nil
}

func startHTTPServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := ggrpc.Dial(*grpcServerEndpoint, ggrpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	opts := []ggrpc.DialOption{ggrpc.WithInsecure()}
	err = grpc.RegisterUsersHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	log.Println("Running HTTP Server at " + *httpServerEndpoint)
	err = http.ListenAndServe(*httpServerEndpoint, mux)
	if err != nil {
		return err
	}
	return nil
}
