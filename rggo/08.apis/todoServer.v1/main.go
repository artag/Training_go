package main

import (
	"flag"     // To handle command-line options
	"fmt"      // To format output
	"net/http" // To handle HTTP connections
	"os"       // For operating system-related functions
	"time"     // To define variables based on time to handle timeouts
)

func main() {
	host := flag.String("h", "localhost", "Server host")
	port := flag.Int("p", 8080, "Server port")
	todoFile := flag.String("f", "todoServer.json", "todo JSON file")
	flag.Parse()

	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *host, *port),
		Handler:      newMux(*todoFile),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
