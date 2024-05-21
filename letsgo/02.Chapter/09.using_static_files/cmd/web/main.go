// Disabling directory listings.
// Add blank index.html to every directory in ./ui/static:
// find ./ui/static -type d -exec touch {}/index.html \;

// To test:
// curl -i -X GET http://localhost:4000/snippet/view?id=123
// curl -i -X POST http://localhost:4000/snippet/create
// http://localhost:4000/static/

package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Create a file server for static files.
	fileServer := http.FileServer(http.Dir("./ui/static/")) // path relative to the root dir.

	// Register the file server as the handler for all URL paths that
	// starts with "/static/".
	// For matching paths, we strip the "/static/" prefix before
	// the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes as normal.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
