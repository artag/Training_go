package main

import (
	"os/exec" // To execute external programs
)

type step struct {
	// Step Name.
	Name string
	// Executable name of the external tool.
	Exe string
	// Arguments for the executable.
	Args []string
	// Output Message in case of success.
	Message string
	// Target project on which to execute the task.
	Proj string
}

func newStep(name, exe, message, proj string, args []string) step {
	return step{
		Name:    name,
		Exe:     exe,
		Message: message,
		Args:    args,
		Proj:    proj,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.Exe, s.Args...)
	// Specify the working directory for the command
	cmd.Dir = s.Proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.Name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.Message, nil
}
