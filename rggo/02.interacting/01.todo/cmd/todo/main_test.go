package main_test

import (
	"fmt"           // To print formatted output
	"os"            // To use operating system types
	"os/exec"       // To execute external commands
	"path/filepath" // To deal with directory paths
	"runtime"       // To identify the running operating system
	"strings"       // To compare strings
	"testing"       // To access testing tools
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

// Call build tool.
// Execute the tests.
// Clean up produced files.
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

// Execute tests that depend on each other
func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	// Subtest 1
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, strings.Split(task, " ")...)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Subtest 2
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := task + "\n"
		if expected != string(out) {
			t.Errorf("Expected '%q', got '%q' instead\n", expected, string(out))
		}
	})
}
