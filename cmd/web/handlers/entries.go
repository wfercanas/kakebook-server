package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
	"github.com/wfercanas/kakebook-server/internal/model"
)

func GetEntryById(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		entryId, err := uuid.Parse(r.PathValue("entryId"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("Invalid Entry Id: %s", r.PathValue("entryId")))
			return
		}

		entry, err := app.Entries.Get(entryId)
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

		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	}
}

func CreateNewEntry(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newEntry model.NewEntry
		err := json.NewDecoder(r.Body).Decode(&newEntry)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		debits := 0.0
		credits := 0.0
		var accountIds []uuid.UUID

		for _, movement := range newEntry.Movements {
			if slices.Contains(accountIds, movement.AccountId) {
				app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("An account can only be used once: %s", movement.AccountId))
				return
			} else {
				accountIds = append(accountIds, movement.AccountId)
			}
		}

		for _, movement := range newEntry.Movements {
			if movement.MovementType == "debit" {
				debits += movement.Value
				continue
			}
			if movement.MovementType == "credit" {
				credits += movement.Value
				continue
			}
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("Invalid Movement Type: %s", movement.MovementType))
			return
		}

		if debits != credits {
			app.ClientError(w, r, http.StatusBadRequest, "Debits and Credits must be equal")
			return
		}

		newEntry.Amount = debits

		err = app.Entries.Insert(newEntry)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteEntry(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		entryId, err := uuid.Parse(r.PathValue("entryId"))
		if err != nil {
			app.ClientError(w, r, http.StatusBadRequest, fmt.Sprintf("Invalid Entry Id: %s", r.PathValue("entryId")))
			return
		}

		err = app.Entries.Delete(entryId)
		if err != nil {
			if errors.Is(err, model.ErrNoRecord) {
				app.ClientError(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			}
			app.ServerError(w, r, err)
		}

		app.Logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusNoContent)
	}
}
