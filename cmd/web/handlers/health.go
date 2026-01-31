package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/wfercanas/kakebook-server/cmd/web/config"
)

func Health(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")
		app.Logger.Info(http.StatusText(http.StatusOK), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		fmt.Fprintf(w, "Hello from kakebook")
	}
}
