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
	Entries  *model.EntryModel
	Accounts *model.AccountModel
	Projects *model.ProjectModel
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.Logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, r *http.Request, status int) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.Logger.Info(http.StatusText(status), slog.String("status", strconv.Itoa(status)), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(status), status)
}
