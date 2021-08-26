package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/pb/v1"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	grpcServerAddr = flag.String("grpc-server-endpoint", "localhost:8070", "gRPC Server Address")
	httpServerAddr = flag.String("http-server-endpoint", "localhost:8080", "HTTP Server Address")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s/n", err)
		os.Exit(1)
	}
}

func run() error {
	err := startHTTPServer()
	if err != nil {
		return errors.Wrap(err, "Start HTTP Server")
	}
	return nil
}

func startHTTPServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(*grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterUserHandlerFromEndpoint(ctx, mux, *grpcServerAddr, opts)
	if err != nil {
		return err
	}

	log.Println("Running HTTP Server at " + *httpServerAddr)
	err = http.ListenAndServe(*httpServerAddr, mux)
	if err != nil {
		return err
	}
	return nil
}
