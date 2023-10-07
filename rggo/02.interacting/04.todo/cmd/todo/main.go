package main

import (
	"bufio"   // Read data from STDIN input stream (os.Stdin)
	"flag"    // To process input
	"fmt"     // To process output
	"io"      // To use io.Reader interface
	"os"      // To verify the arguments from cli
	"strings" // To use Join() to compose a task name
	"time"

	"rggo/interacting/todo"
)

const (
	environmentVariable = "TODO_FILENAME"
)

// Default file name
var todoFileName = ".todo.json"

func main() {
	checkEnvironment()
	redefineFlagUsage()

	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	l := &todo.List{}

	// Read todo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		// List current todo items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Save the new list
		saveListToFile(l)
	case *add:
		// Add the task.
		// When any arguments (excluding flags) are provided, they will be used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		// Save the new list
		saveListToFile(l)
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func checkEnvironment() {
	if os.Getenv(environmentVariable) != "" {
		todoFileName = os.Getenv(environmentVariable)
	}
}

func redefineFlagUsage() {
	appName := strings.Trim(os.Args[0], "./")
	year := time.Now().Year()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Copyright %d.\n", appName, year)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}
}

// Save the list to the file
func saveListToFile(l *todo.List) {
	if err := l.Save(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Get description for a new task. From arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}

	return s.Text(), nil
}
