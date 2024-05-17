package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Check if current request URL path exactly matches '/'. If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

// Add a snippetCreate hanlder function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet"))
}

func main() {
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)                        // subtree path example (ends with '/')
	mux.HandleFunc("/snippet/view", snippetView)     // fixed path example
	mux.HandleFunc("/snippet/create", snippetCreate) // fixed path example

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
