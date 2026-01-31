package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

type NewEntry struct {
	ProjectId   uuid.UUID `json:"project_id"`
	Date        string    `json:"date"`
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
}

func (app *application) createNewEntry(w http.ResponseWriter, r *http.Request) {
	var entry NewEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	date, err := time.Parse("2006-01-02", entry.Date)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.entries.Insert(date, entry.Description, entry.ProjectId, entry.Amount)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	w.WriteHeader(http.StatusCreated)
}

type NewAccount struct {
	Name            string    `json:"name"`
	AccountCategory string    `json:"account_category"`
	ProjectId       uuid.UUID `json:"project_id"`
}

func (app *application) createNewAccount(w http.ResponseWriter, r *http.Request) {
	var account NewAccount
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if account.Name == "" || account.AccountCategory == "" {
		app.logger.Info("Missing name or account_category")
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	if account.ProjectId == uuid.Nil {
		app.logger.Info("Missing project_id")
		app.clientError(w, r, http.StatusBadRequest)
		return
	}

	err = app.accounts.Insert(account.Name, account.AccountCategory, account.ProjectId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	w.WriteHeader(http.StatusCreated)
}
