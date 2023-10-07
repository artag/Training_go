package main

import (
	"flag" // To process input
	"fmt"  // To process output
	"os"   // To verify the arguments from cli
	"strings"
	"time"

	"rggo/interacting/todo"
)

// Hardcoding the file name
const todoFileName = ".todo.json"

func main() {
	redefineFlagUsage()

	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the todo list")
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
	case *task != "":
		// Add the task
		l.Add(*task)
		// Save the new list
		saveListToFile(l)
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
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
