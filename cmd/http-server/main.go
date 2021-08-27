package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/internal/pb/v1"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	grpcServerAddr string
	httpServerAddr string
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s/n", err)
		os.Exit(1)
	}
}

func run() error {
	grpcServerAddr = os.Getenv("GRPC_SERVER_ADDR")
	httpServerAddr = os.Getenv("HTTP_SERVER_ADDR")

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

	conn, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, rmux, grpcServerAddr, opts)
	if err != nil {
		return err
	}

	// Swagger
	smux := http.NewServeMux()
	smux.Handle("/", rmux)
	smux.HandleFunc("/swagger.json", handleSwagger)
	fs := http.FileServer(http.Dir("www/swagger-ui"))
	smux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	log.Println("Running HTTP Server at " + httpServerAddr)
	log.Println("Swagger is at: " + httpServerAddr + "/swagger-ui/")
	err = http.ListenAndServe(httpServerAddr, smux)
	if err != nil {
		return err
	}
	return nil
}

func handleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/swagger.json")
}
