package pomodoro_test

import (
	"rggo/interactiveTools/pomo/pomodoro"
	"rggo/interactiveTools/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
