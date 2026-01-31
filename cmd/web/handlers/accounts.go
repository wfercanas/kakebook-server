package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
)

type newAccount struct {
	Name            string    `json:"name"`
	AccountCategory string    `json:"account_category"`
	ProjectId       uuid.UUID `json:"project_id"`
}

func CreateNewAccount(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var account newAccount
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		if account.Name == "" || account.AccountCategory == "" {
			app.Logger.Info("Missing name or account_category")
			app.ClientError(w, r, http.StatusBadRequest)
			return
		}

		if account.ProjectId == uuid.Nil {
			app.Logger.Info("Missing project_id")
			app.ClientError(w, r, http.StatusBadRequest)
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
