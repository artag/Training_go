package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Represents a todo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// Represents a list of todo items
type List []item

// Creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Marks a todo item as completed
// by setting Done = true and CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exists", i)
	}

	idx := i - 1
	ls[idx].Done = true
	ls[idx].CompletedAt = time.Now()

	return nil
}

// Deletes a todo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exists", i)
	}

	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

// Encodes the List as JSON and saves it using the provided file name
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Opens the provided file name, decodes the JSON data and parses it into a List
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

// Prints out a formatted list.
func (l *List) Print(verbose bool, exludeCompleted bool) string {
	formatted := ""

	for k, t := range *l {
		prefix := "  "
		if t.Done {
			if exludeCompleted {
				continue
			}

			prefix = "X "
		}

		dateCreated, dateCompleted := l.getDatesAsString(verbose, &t)

		// Adjust the item number k to print numbers starting from 1 instead of 0
		formatted += fmt.Sprintf("%s%d: %s    %s    %s\n", prefix, k+1, t.Task, dateCreated, dateCompleted)
	}

	return formatted
}

func (*List) getDatesAsString(verbose bool, t *item) (string, string) {
	dateCreated := ""
	dateCompleted := ""
	if !verbose {
		return dateCreated, dateCompleted
	}

	fmtDateCreated := t.CreatedAt.Format("2006-01-02 15:04:05")
	dateCreated = fmt.Sprintf("Created at: %s", fmtDateCreated)

	if t.Done {
		fmtDateCompleted := t.CompletedAt.Format("2006-01-02 15:04:05")
		dateCompleted = fmt.Sprintf("Completed at: %s", fmtDateCompleted)
	}

	return dateCreated, dateCompleted
}
