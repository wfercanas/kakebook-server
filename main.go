package main

import (
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from kakebook"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
