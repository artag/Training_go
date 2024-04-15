//go:build inmemory || containers
// +build inmemory containers

package cmd

import (
	"rggo/interactiveTools/pomo/pomodoro"
	"rggo/interactiveTools/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
