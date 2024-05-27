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

type config struct {
	addr      string
	staticDir string
}

func main() {
	// Get variable from environment
	// addr := os.Getenv("SNIPPETBOX_ADDR")

	// Or (better): Define a command line flag.
	// Store all configuration settings in a single struct.
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	// Parse the command line flags. Need to call before using flags.
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use 'addr' flag.
	// Ports 0-1023 using by services with root privileges.
	log.Printf("Starting server on %s", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
}
