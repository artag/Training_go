package main

import (
	"fmt"     // To process output
	"os"      // To verify the arguments from cli
	"strings" // To process input

	"rggo/interacting/todo"
)

// Hardcoding the file name
const todoFileName = ".todo.json"

func main() {
	l := &todo.List{}

	// Read todo items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	// Without extra arguments, print the list
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	// Concatenate all provided arguments with a space and add to the list as an item
	default:
		item := strings.Join(os.Args[1:], " ")
		// Add the task
		l.Add(item)

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
