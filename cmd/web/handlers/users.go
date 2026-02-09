package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
	"github.com/wfercanas/kakebook-server/internal/model"
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
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound)
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		fmt.Fprintf(w, user.String())
	}
}

func GetProjectsByUserId(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := uuid.Parse(r.PathValue("userID"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest)
			return
		}

		projects, err := app.Projects.GetProjectsByUserId(userId)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound)
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		var body struct {
			Projects []model.Project `json:"projects"`
		}
		body.Projects = append(body.Projects, projects...)

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(body)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	}
}

func CreateNewUser(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Created a new user...")
	}
}
