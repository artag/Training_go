//go:build inmemory
// +build inmemory

package cmd

import (
	"rggo/interactiveTools/pomo/pomodoro"
	"rggo/interactiveTools/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}