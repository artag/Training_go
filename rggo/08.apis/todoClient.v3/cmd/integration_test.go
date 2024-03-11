//go:build integration
// +build integration

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func randomTaskName(t *testing.T) string {
	t.Helper()
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var p strings.Builder
	for i := 0; i < 32; i++ {
		p.WriteByte(chars[r.Intn(len(chars))])
	}

	return p.String()
}

func TestIntegration(t *testing.T) {
	apiRoot := "http://localhost:8080"
	if os.Getenv("TODO_API_ROOT") != "" {
		apiRoot = os.Getenv("TODO_API_ROOT")
	}

	today := time.Now().Format("Jan/02")
	task := randomTaskName(t)
	taskId := ""

	// Step 1.
	t.Run("1.AddTask", func(t *testing.T) {
		args := []string{task}
		expOut := fmt.Sprintf("Added task %q to the list.\n", task)

		// Execute Add test
		var out bytes.Buffer
		if err := addAction(&out, apiRoot, args); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}
		if expOut != out.String() {
			t.Errorf("Expected output %q, got %q.", expOut, out.String())
		}
	})

	// Step 2.
	t.Run("2.ListTasks", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}

		// Execute List test
		outList := ""
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outList = scanner.Text()
				break
			}
		}
		if outList == "" {
			t.Errorf("Task %q is not in the list", task)
		}

		taskCompleteStatus := strings.Fields(outList)[0]
		if taskCompleteStatus != "-" {
			t.Errorf("Expected status %q, got %q.", "-", taskCompleteStatus)
		}

		taskId = strings.Fields(outList)[1]
	})

	// Step 3.
	vRes := t.Run("3.ViewTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := viewAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}

		// Execute View test
		viewOut := strings.Split(out.String(), "\n")

		if !strings.Contains(viewOut[0], task) {
			t.Fatalf("Expected task %q, got %q.", task, viewOut[0])
		}
		if !strings.Contains(viewOut[1], today) {
			t.Fatalf("Expected creation day/month %q, got %q.", today, viewOut[1])
		}
		if !strings.Contains(viewOut[2], "No") {
			t.Fatalf("Expected completed status %q, got %q.", "No", viewOut[2])
		}
	})
	if !vRes {
		t.Fatalf("View task failed. Stopping integration tests.")
	}

	// Step 4.
	t.Run("4.CompleteTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := completeAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}

		expOut := fmt.Sprintf("Item number %s marked as completed.\n", taskId)
		if expOut != out.String() {
			t.Fatalf("Expected output %q, got %q.", expOut, out.String())
		}
	})

	// Step 5.
	t.Run("5.ListCompletedTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}

		outlist := ""
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), task) {
				outlist = scanner.Text()
				break
			}
		}

		if outlist == "" {
			t.Errorf("Task %q is not in the list", task)
		}

		taskCompleteStatus := strings.Fields(outlist)[0]
		if taskCompleteStatus != "X" {
			t.Errorf("Expected status %q, got %q.", "X", taskCompleteStatus)
		}
	})

	// Step 6.
	t.Run("6.DeleteTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := delAction(&out, apiRoot, taskId); err != nil {
			t.Fatalf("Expected no error, got %q.", err)
		}

		expOut := fmt.Sprintf("Item number %s deleted.\n", taskId)
		if expOut != out.String() {
			t.Fatalf("Expected output %q, got %q.", expOut, out.String())
		}
	})

	// Step 7.
	t.Run("7.ListDeletedTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			expErr := "Not found: No results found"
			if err.Error() != expErr {
				t.Fatalf("Expected error %q, got %q.", expErr, err.Error())
			}
		} else {
			t.Fatal("Expected error.")
		}
	})
}
