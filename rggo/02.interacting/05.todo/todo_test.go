package todo_test

import (
	"fmt"
	"io/ioutil" // to create temprary files
	"os"        // to delete temporary files
	"rggo/interacting/todo"
	"testing"
	"time"
)

// Tests the Add metod of the List type
func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "New Task"

	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l[0].Task)
	}
}

// Tests the Complete method of the List type
func TestComplete(t *testing.T) {
	// Arrange
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New task not be completed.")
	}

	// Act
	l.Complete(1)

	// Assert
	if !l[0].Done {
		t.Errorf("New task should be completed.")
	}
}

// Tests the Delete method of the List type
func TestDelete(t *testing.T) {
	// Arrange
	l := todo.List{}
	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}
	for _, t := range tasks {
		l.Add(t)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Expected %q, got %q instead.", tasks[0], l[0].Task)
	}

	// Act
	l.Delete(2)

	// Assert
	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead.", tasks[2], l[1].Task)
	}
}

// Tests the Save and Get methods of the List type
func TestSaveGet(t *testing.T) {
	// Arrange
	l1 := todo.List{}
	l2 := todo.List{}
	taskName := "New Task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l1[0].Task)
	}

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	// Act, Assert
	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	// Act, Assert
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Fatalf("Task %q should match %q task.", l1[0].Task, l2[0].Task)
	}
}

// Test print all list items to string without verbose
func TestPrintList(t *testing.T) {
	// Arrange
	l := todo.List{}
	l.Add("Task 1")
	l.Add("Task 2")
	l.Add("Task 3")
	l.Complete(2)

	expected :=
		"  1: Task 1        \n" +
			"X 2: Task 2        \n" +
			"  3: Task 3        \n"

	// Act
	actual := l.Print(false, false)

	// Assert
	assertString(expected, actual, t)
}

// Test verbose print all list items to string
func TestVerbosePrintList(t *testing.T) {
	// Arrange
	l := todo.List{}
	l.Add("Task 1")
	dt1 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second)

	l.Add("Task 2")
	dt2 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second)

	l.Add("Task 3")
	dt3 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second)

	l.Complete(2)
	dt4 := time.Now().Format("2006-01-02 15:04:05")

	expected := fmt.Sprintf(
		"  1: Task 1    Created at: %s    \n"+
			"X 2: Task 2    Created at: %s    Completed at: %s\n"+
			"  3: Task 3    Created at: %s    \n",
		dt1, dt2, dt4, dt3)

	// Act
	actual := l.Print(true, false)

	// Assert
	assertString(expected, actual, t)
}

// Test print list without completed items to string
func TestPrintListWithoutCompletedTasks(t *testing.T) {
	// Arrange
	l := todo.List{}
	l.Add("Task 1")
	l.Add("Task 2")
	l.Add("Task 3")
	l.Complete(2)

	expected :=
		"  1: Task 1        \n" +
			"  3: Task 3        \n"

	// Act
	actual := l.Print(false, true)

	// Assert
	assertString(expected, actual, t)
}

// Test verbose print list without completed items to string
func TestVerbosePrintListWithoutCompletedTasks(t *testing.T) {
	// Arrange
	l := todo.List{}
	l.Add("Task 1")
	dt1 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second)

	l.Add("Task 2")
	time.Sleep(time.Second)

	l.Add("Task 3")
	dt3 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second)

	l.Complete(2)

	expected := fmt.Sprintf(
		"  1: Task 1    Created at: %s    \n"+
			"  3: Task 3    Created at: %s    \n",
		dt1, dt3)

	// Act
	actual := l.Print(true, true)

	// Assert
	assertString(expected, actual, t)
}

func assertString(expected string, actual string, t *testing.T) {
	if expected == actual {
		return
	}

	t.Fatalf(
		"Expected:\n"+
			"%q\n"+
			"Actual:\n"+
			"%q\n",
		expected, actual)
}
