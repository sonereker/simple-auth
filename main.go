package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sonereker/simple-auth/internal"
	"github.com/sonereker/simple-auth/users"
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
	db, err := internal.NewDBConnection()
	if err != nil {
		return errors.Wrap(err, "Init database")
	}

	err = db.AutoMigrate(&users.UserDBModel{})
	if err != nil {
		return errors.Wrap(err, "Run migrations")
	}

	router := internal.InitRouter()
	server := internal.NewServer(db, router)

	userHandler := users.NewHandler(server)
	userHandler.RegisterRoutes()

	fmt.Printf("Running Simple Auth Server")
	err = http.ListenAndServe(":8080", server.Router)
	if err != nil {
		return errors.Wrap(err, "Run API server")
	}
	return nil
}
