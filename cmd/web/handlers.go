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

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	fmt.Fprint(w, "Returning all users...")
}

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	user, err := app.users.Get(userId)
	if err != nil {
		app.clientError(w, r, http.StatusNotFound)
		app.logger.Info(err.Error())
		return
	}

	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	fmt.Fprintf(w, user.String())
}

func (app *application) createNewUser(w http.ResponseWriter, r *http.Request) {
	app.logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created a new user...")
}
