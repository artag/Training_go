// To test:
// curl -i -X POST http://localhost:4000/snippet/create

package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// Check if current request URL path exactly matches '/'.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {
		w.WriteHeader(405) // 405 - Mehod not allowed
		w.Write([]byte("Method Not Allowed"))
		return
	}

	w.Write([]byte("Create a new snippet...")) // 200 - OK
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)                        // subtree path example (ends with '/')
	mux.HandleFunc("/snippet/view", snippetView)     // fixed path example
	mux.HandleFunc("/snippet/create", snippetCreate) // fixed path example

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
