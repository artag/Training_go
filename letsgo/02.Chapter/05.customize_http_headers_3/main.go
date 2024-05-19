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
	if r.Method != http.MethodPost {
		// Set a new-cache-control header. Overwrite existing value.
		// w.Header().Set("Cache-Control", "public, max-age=31536000")

		// The Add() method appends header and can be called multiple times.
		// w.Header().Add("Cache-Control", "public")
		// w.Header().Add("Cache-Control", "max-age=31536000")

		// Delete all values for the "Cache-Control" header.
		// w.Header().Del("Cache-Control")

		// Retrieve the first value for the "Cache-Control" header.
		// w.Header().Get("Cache-Control")

		// Retrieves a slice of all values for the "Cache-Control" header.
		// w.Header().Values("Cache-Control")

		w.Header().Set("Allow", http.MethodPost)
		// Use the http.Error() function to send a 405 status code and
		// "Method Not Allowed" string as the response body.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // 405
		return
	}

	w.Write([]byte("Create a new snippet...")) // 200 - OK
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
