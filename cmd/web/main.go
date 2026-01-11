package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Address")
	flag.Parse()

	app := application{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /users", users)
	mux.HandleFunc("GET /users/{userID}", getUserById)

	mux.HandleFunc("POST /users", createNewUser)

	app.logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
