package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, status int) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Info(http.StatusText(status), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, http.StatusText(status), status)
}
