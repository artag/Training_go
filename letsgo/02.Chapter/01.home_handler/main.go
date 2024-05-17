package main

import (
	"log"
	"net/http"
)

// Define handler function which writes a byte slice containing
// string as the reponse body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {
	// Use the http.NewServerMux() to initialize a new sevemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Use the http.ListenAndServe() function to start a new web server.
	// Format of server address is "host:port".
	// If omit "host" than server will listen on all computers available
	// network interfaces.
	// If using named ports like ":http" or ":http-alt" then Go look up
	// port number from /etc/services.
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
