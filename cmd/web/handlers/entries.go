package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wfercanas/kakebook-server/cmd/web/config"
)

type newEntry struct {
	ProjectId   uuid.UUID `json:"project_id"`
	Date        string    `json:"date"`
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
}

func CreateNewEntry(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var entry newEntry
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		date, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		err = app.Entries.Insert(date, entry.Description, entry.ProjectId, entry.Amount)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		app.Logger.Info(http.StatusText(http.StatusCreated), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		w.WriteHeader(http.StatusCreated)
	}

}
