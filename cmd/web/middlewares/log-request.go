package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/wfercanas/kakebook-server/cmd/web/config"
)

func LogRequest(app *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("request received", slog.String("ip", r.RemoteAddr), slog.String("proto", r.Proto), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}
