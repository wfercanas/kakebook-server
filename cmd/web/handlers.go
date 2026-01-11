package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	fmt.Fprintf(w, "Hello from kakebook")
}

func users(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Returning all users...")
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Returning the specific user with ID: %s", id.String())
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Created a new user...")
}
