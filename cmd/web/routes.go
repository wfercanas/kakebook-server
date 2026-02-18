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
	mux.HandleFunc("GET /users/{userId}", handlers.GetUserById(app))
	mux.HandleFunc("GET /users/{userId}/projects", handlers.GetProjectsByUserId(app))
	mux.HandleFunc("GET /projects/{projectId}/accounts", handlers.GetAccountsByProjectId(app))
	mux.HandleFunc("GET /accounts/{accountId}", handlers.GetAccountById(app))
	mux.HandleFunc("GET /accounts/{accountId}/ledger", handlers.GetAccountLedgerById(app))
	mux.HandleFunc("GET /entries/{entryId}", handlers.GetEntryById(app))

	mux.HandleFunc("POST /users", handlers.CreateNewUser(app))
	mux.HandleFunc("POST /entries", handlers.CreateNewEntry(app))
	mux.HandleFunc("POST /accounts", handlers.CreateNewAccount(app))

	mux.HandleFunc("DELETE /entries/{entryId}", handlers.DeleteEntry(app))

	return mux
}
