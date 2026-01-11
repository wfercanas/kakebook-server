package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /users", users)
	mux.HandleFunc("GET /users/{userID}", getUserById)

	mux.HandleFunc("POST /users", createNewUser)

	log.Print("Starting server listening on ", *addr)

	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
