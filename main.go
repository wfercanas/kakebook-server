package main

import (
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from kakebook"))
}

func users(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Returning all users..."))
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Returning the specific user with ID: " + id.String()))
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created a new user..."))
}

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
