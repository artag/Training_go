// Redirect the stdout and stderr streams to on-disk files on starting:
// go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Create a logger for writing information messages to standard output.
	// With local date and local time.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages to standard error output.
	// Lshortfile - include relevant filename and line number.
	// Llongfile - include the full filepath.
	// LUTC - to use UTC datetimes (instead of local ones).
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Initialize a new http server using the custom errorLog logger.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Write messages using the two new loggers.
	infoLog.Printf("Starting server on %s", *addr)
	// Call the ListenAndServe() method on http.Server struct.
	err := srv.ListenAndServe()
	errorLog.Fatal(err) // Recommended panic or exit directly from main() only.
}
