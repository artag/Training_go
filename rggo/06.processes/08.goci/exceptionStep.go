package main

import (
	"bytes"   // To capture the program's output
	"fmt"     // To format message
	"os/exec" // To execute external programs
)

type exceptionStep struct {
	step
}

func newExceptionStep(name, exe, message, proj string, args []string) exceptionStep {
	s := exceptionStep{}
	s.step = newStep(name, exe, message, proj, args)
	return s
}

func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.Exe, s.Args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = s.Proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.Name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepErr{
			step:  s.Name,
			msg:   fmt.Sprintf("invalid format: %s", out.String()),
			cause: nil,
		}
	}

	return s.Message, nil
}
