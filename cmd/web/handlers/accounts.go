package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
	"github.com/wfercanas/kakebook-server/internal/model"
)

func GetAccountById(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := uuid.Parse(r.PathValue("accountID"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, "Invalid Account Id")
			return
		}

		account, err := app.Accounts.GetAccountById(accountId)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(account)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	}
}

func CreateNewAccount(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var account model.NewAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		if account.Name == "" {
			app.ClientError(w, r, http.StatusBadRequest, "Invalid Account Name")
			return
		}

		if account.AccountCategory == "" {
			app.ClientError(w, r, http.StatusBadRequest, "Invalid Account Category")
			return
		}

		if account.ProjectId == uuid.Nil {
			app.ClientError(w, r, http.StatusBadRequest, "Invalid Project Id")
			return
		}

		err = app.Accounts.Insert(account.Name, account.AccountCategory, account.ProjectId)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusCreated)
	}
}
