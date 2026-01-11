package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /users", users)
	mux.HandleFunc("GET /users/{userID}", getUserById)

	mux.HandleFunc("POST /users", createNewUser)

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
