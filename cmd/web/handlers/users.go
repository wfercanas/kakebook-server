package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
)

func GetUsers(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		fmt.Fprint(w, "Returning all users...")
	}
}

func GetUserById(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := uuid.Parse(r.PathValue("userID"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest)
			return
		}

		user, err := app.Users.Get(userId)
		if err != nil {
			app.ClientError(w, r, http.StatusNotFound)
			app.Logger.Info(err.Error())
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		fmt.Fprintf(w, user.String())
	}
}

func CreateNewUser(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Created a new user...")
	}
}
