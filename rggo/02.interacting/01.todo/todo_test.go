package todo_test

import (
	"io/ioutil" // to create temprary files
	"os"        // to delete temporary files
	"rggo/interacting/todo"
	"testing"
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
