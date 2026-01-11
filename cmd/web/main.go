package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /users", users)
	mux.HandleFunc("GET /users/{userID}", getUserById)

	mux.HandleFunc("POST /users", createNewUser)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
