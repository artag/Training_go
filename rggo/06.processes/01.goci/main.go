package main

import (
	"flag"    // To parse command line options
	"fmt"     // To handle output
	"io"      // Provides io.Writer interface
	"os"      // To interact with the operating system
	"os/exec" // To execute external programs
)

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}

	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)

	// Specify the working directory for the command
	cmd.Dir = proj
	if err := cmd.Run(); err != nil {
		return &stepErr{step: "go build", msg: "go build failed", cause: err}
	}

	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	return err
}
