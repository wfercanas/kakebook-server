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

func GetAccountLedgerById(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := uuid.Parse(r.PathValue("accountId"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, "Invalid Account Id")
			return
		}

		account, err := app.Ledger.GetLedgerAccountById(accountId)
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
