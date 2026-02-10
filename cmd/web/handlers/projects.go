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

func GetAccountsByProjectId(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, err := uuid.Parse(r.PathValue("projectID"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest)
			return
		}

		accounts, err := app.Projects.GetAccountsByProjectId(projectId)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound)
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		for i := range accounts {
			app.Accounts.CalculateAccountBalance(&accounts[i])
		}

		var body struct {
			Accounts []model.Account `json:"accounts"`
		}
		body.Accounts = accounts

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(body)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	}
}
