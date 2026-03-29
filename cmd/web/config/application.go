package config

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/wfercanas/kakebook-server/internal/model"
)

type Application struct {
	Logger   *slog.Logger
	Users    *model.UserModel
	Journal  *model.JournalModel
	Accounts *model.AccountModel
	Projects *model.ProjectModel
	Ledger   *model.LedgerModel
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.Logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.Logger.Info(http.StatusText(status), slog.String("status", strconv.Itoa(status)), slog.String("method", method), slog.String("uri", uri), slog.String("msg", message))
	http.Error(w, message, status)
}

func (app *Application) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("request received", slog.String("ip", r.RemoteAddr), slog.String("proto", r.Proto), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}
