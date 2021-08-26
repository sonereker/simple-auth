package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/pb"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	grpcServerAddr = flag.String("grpc-server-endpoint", ":8070", "gRPC Server Address")
	httpServerAddr = flag.String("http-server-endpoint", ":8080", "HTTP Server Address")
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

	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, rmux, *grpcServerAddr, opts)
	if err != nil {
		return err
	}

	// Swagger
	smux := http.NewServeMux()
	smux.Handle("/", rmux)
	smux.HandleFunc("/swagger.json", handleSwagger)
	fs := http.FileServer(http.Dir("www/swagger-ui"))
	smux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	log.Println("Running HTTP Server at " + *httpServerAddr)
	log.Println("Swagger is at: http://localhost:8080/swagger-ui/")
	err = http.ListenAndServe(*httpServerAddr, smux)
	if err != nil {
		return err
	}
	return nil
}

func handleSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/swagger.json")
}
