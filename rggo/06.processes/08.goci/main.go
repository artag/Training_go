package main

import (
	"flag"      // To parse command line options
	"fmt"       // To handle output
	"io"        // Provides io.Writer interface
	"os"        // To interact with the operating system
	"os/signal" // To handle signals
	"syscall"   // To use signal definitions
)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	// Save2()
	// return

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}

	// Load pipeline from config file.
	pipeline, err := LoadPipeline(proj)
	if err != nil {
		return err
	}

	// Signal definition. We are interested in two terminating signals only.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Error channel to communicate potential errors.
	errCh := make(chan error)
	// To communicate the loop conclusion.
	done := make(chan struct{})

	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}

			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()

	for {
		select {

		// Handle the signal.
		case rec := <-sig:
			// Stop receiving more signals on the sig channel.
			signal.Stop(sig)
			return fmt.Errorf("%s: Exiting: %w", rec, ErrSignal)

		// Handle error.
		case err := <-errCh:
			return err

			// Handle success finish
		case <-done:
			return nil
		}
	}
}
