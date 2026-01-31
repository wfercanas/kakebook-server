package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /users", app.getUsers)
	mux.HandleFunc("GET /users/{userID}", app.getUserById)

	mux.HandleFunc("POST /users", app.createNewUser)
	mux.HandleFunc("POST /entries", app.createNewEntry)
	mux.HandleFunc("POST /accounts", app.createNewAccount)

	return mux
}
