package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	fmt.Fprintf(w, "Hello from kakebook")
}

func (app *application) users(w http.ResponseWriter, r *http.Request) {
	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	fmt.Fprint(w, "Returning all users...")
}

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		app.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		http.NotFound(w, r)
		return
	}

	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	fmt.Fprintf(w, "Returning the specific user with ID: %s", id.String())
}

func (app *application) createNewUser(w http.ResponseWriter, r *http.Request) {
	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created a new user...")
}
