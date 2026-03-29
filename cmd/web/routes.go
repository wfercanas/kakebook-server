package main

import (
	"net/http"

	"github.com/wfercanas/kakebook-server/cmd/web/config"
	"github.com/wfercanas/kakebook-server/cmd/web/handlers"
	"github.com/wfercanas/kakebook-server/cmd/web/middlewares"
)

func routes(app *config.Application) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/{$}", handlers.Health(app))
	mux.HandleFunc("GET /api/users/{$}", handlers.GetUsers(app))
	mux.HandleFunc("GET /api/users/{userId}", handlers.GetUserById(app))
	mux.HandleFunc("GET /api/users/{userId}/projects", handlers.GetProjectsByUserId(app))
	mux.HandleFunc("GET /api/projects/{projectId}/accounts", handlers.GetAccountsByProjectId(app))
	mux.HandleFunc("GET /api/accounts/{accountId}", handlers.GetAccountById(app))
	mux.HandleFunc("GET /api/accounts/{accountId}/ledger", handlers.GetAccountLedgerById(app))
	mux.HandleFunc("GET /api/projects/{projectId}/journal", handlers.GetJournalByProjectId(app))
	mux.HandleFunc("GET /api/entries/{entryId}", handlers.GetEntryById(app))

	mux.HandleFunc("POST /api/users", handlers.CreateNewUser(app))
	mux.HandleFunc("POST /api/entries", handlers.CreateNewEntry(app))
	mux.HandleFunc("POST /api/accounts", handlers.CreateNewAccount(app))

	mux.HandleFunc("DELETE /api/accounts/{accountId}", handlers.DeleteAccount(app))
	mux.HandleFunc("DELETE /api/entries/{entryId}", handlers.DeleteEntry(app))

	fs := http.FileServer(http.Dir("/ui/dist"))
	mux.HandleFunc("/", handlers.Frontend(fs, "/ui/dist/index.html"))

	return middlewares.LogRequest(app, mux)
}
