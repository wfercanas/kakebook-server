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

func GetAccountsByProjectId(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, err := uuid.Parse(r.PathValue("projectId"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("Invalid Project Id: %s", r.PathValue("projectId")))
			return
		}

		accounts, err := app.Projects.GetAccountsByProjectId(projectId)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				app.ServerError(w, r, err)
				return
			}
		}

		groups := make(map[string][]model.Account)
		for i, account := range accounts {
			err = app.Accounts.CalculateAccountBalance(&accounts[i])
			if err != nil {
				app.ServerError(w, r, err)
				return
			}
			groups[account.AccountCategory] = append(groups[account.AccountCategory], accounts[i])
		}

		type accountsBody struct {
			Assets      []model.Account `json:"assets"`
			Liabilities []model.Account `json:"liabilities"`
			Equity      []model.Account `json:"equity"`
			Revenue     []model.Account `json:"revenue"`
			Expenses    []model.Account `json:"expenses"`
		}

		body := accountsBody{
			Assets:      groups["assets"],
			Liabilities: groups["liabilities"],
			Equity:      groups["equity"],
			Revenue:     groups["revenue"],
			Expenses:    groups["expenses"],
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(body)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	}
}
