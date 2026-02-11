package main

import (
	"net/http"

	"github.com/wfercanas/kakebook-server/cmd/web/config"
	"github.com/wfercanas/kakebook-server/cmd/web/handlers"
)

func routes(app *config.Application) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.Health(app))
	mux.HandleFunc("GET /users/{$}", handlers.GetUsers(app))
	mux.HandleFunc("GET /users/{userID}", handlers.GetUserById(app))
	mux.HandleFunc("GET /users/{userID}/projects", handlers.GetProjectsByUserId(app))
	mux.HandleFunc("GET /projects/{projectID}/accounts", handlers.GetAccountsByProjectId(app))
	mux.HandleFunc("GET /accounts/{accountID}", handlers.GetAccountById(app))

	mux.HandleFunc("POST /users", handlers.CreateNewUser(app))
	mux.HandleFunc("POST /entries", handlers.CreateNewEntry(app))
	mux.HandleFunc("POST /accounts", handlers.CreateNewAccount(app))

	mux.HandleFunc("DELETE /entries/{entryId}", handlers.DeleteEntry(app))

	return mux
}
