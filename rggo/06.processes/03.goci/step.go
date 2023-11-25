package main

import (
	"os/exec" // To execute external programs
)

type step struct {
	// Step name.
	name string
	// Executable name of the external tool.
	exe string
	// Arguments for the executable.
	args []string
	// Output message in case of success.
	message string
	// Target project on which to execute the task.
	proj string
}

func newStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		message: message,
		args:    args,
		proj:    proj,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	// Specify the working directory for the command
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
