//go:build !inmemory
// +build !inmemory

package pomodoro_test

import (
	"os"
	"rggo/interactiveTools/pomo/pomodoro"
	"rggo/interactiveTools/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	tf, err := os.CreateTemp("", "pomo")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	t.Log("Using SQLite repository at " + tf.Name())

	dbRepo, err := repository.NewSQLite3Repo(tf.Name())
	if err != nil {
		t.Fatal(t)
	}

	return dbRepo, func() {
		os.Remove(tf.Name())
	}
}
