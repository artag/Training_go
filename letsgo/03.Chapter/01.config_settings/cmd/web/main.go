// To start server with custom port number:
// go run ./cmd/web -addr=":9999"

// Use environment variables from command line:
// export SNIPPETBOX_ADDR=":9999"
// go run ./cmd/web -addr=$SNIPPETBOX_ADDR

// Show help:
// go run ./cmd/web -help

package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Get variable from environment
	// addr := os.Getenv("SNIPPETBOX_ADDR")

	// Or (better): Define a command line flag.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Parse the command line flags. Need to call before using flags.
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use 'addr' flag.
	// Ports 0-1023 using by services with root privileges.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
