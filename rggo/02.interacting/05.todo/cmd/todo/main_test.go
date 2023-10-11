package main_test

// import . "main"

import (
	"fmt"           // To print formatted output
	"io"            // To use io.WriteString()
	"os"            // To use operating system types
	"os/exec"       // To execute external commands
	"path/filepath" // To deal with directory paths
	"runtime"       // To identify the running operating system
	"strings"
	"testing" // To access testing tools

	globals "rggo/interacting/todo/cmd"
)

var (
	binName = "todo"
)

func TestSplit(t *testing.T) {
	t1 := "task1"
	t2 := "task2"
	txt := fmt.Sprintf("%s\n%s", t1, t2)
	split := strings.Split(txt, "\n")
	assertString("task1", split[0], t)
	assertString("task2", split[1], t)
}

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

	var env string
	if os.Getenv(globals.EnvironmentVariable) != "" {
		env = os.Getenv(globals.EnvironmentVariable)
	}
	os.Setenv(globals.EnvironmentVariable, globals.TestFileName)
	fmt.Printf(
		"Set environment variable '%s=%s' for tests...\n",
		globals.EnvironmentVariable, globals.TestFileName)

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(globals.TestFileName)

	if env != "" {
		os.Setenv(globals.EnvironmentVariable, env)
		fmt.Printf("Restore environment variable '%s=%s' for tests...\n",
			globals.EnvironmentVariable, env)
	}

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
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Subtest 2
	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Subtest 3
	t.Run("ListTasks", func(t *testing.T) {
		out := listCommand(cmdPath, t)

		expected := fmt.Sprintf(
			"  1: %s        \n"+
				"  2: %s        \n",
			task, task2)

		actual := string(out)
		assertString(expected, actual, t)
	})

	// Subtest 4
	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "2")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		out := listCommand(cmdPath, t)
		expected := fmt.Sprintf(
			"  1: %s        \n"+
				"X 2: %s        \n",
			task, task2)

		actual := string(out)
		assertString(expected, actual, t)
	})

	// Subtest 5
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-del", "1")
		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		out := listCommand(cmdPath, t)
		expected := fmt.Sprintf(
			"X 1: %s        \n",
			task2)

		actual := string(out)
		assertString(expected, actual, t)
	})

	// Subtest 6
	task3 := "test task3"
	task4 := "test task4"
	stdin := fmt.Sprintf("%s\n%s", task3, task4)
	t.Run("AddNewTasksFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, stdin)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		out := listCommand(cmdPath, t)
		expected := fmt.Sprintf(
			"X 1: %s        \n"+
				"  2: %s        \n"+
				"  3: %s        \n",
			task2, task3, task4)

		actual := string(out)
		assertString(expected, actual, t)
	})
}

func listCommand(cmdPath string, t *testing.T) []byte {
	cmd := exec.Command(cmdPath, "-list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	return out
}

func assertString(expected string, actual string, t *testing.T) {
	if expected == actual {
		return
	}

	t.Errorf(
		"Expected:\n"+
			"%q\n"+
			"Actual:\n"+
			"%q\n",
		expected, actual)
}
