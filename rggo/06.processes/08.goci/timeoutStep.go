package main

import (
	"context" // To create a context to carry the timeout
	"os/exec" // To execute external commands
	"time"    // To define time values
)

type timeoutStep struct {
	step
	Timeout time.Duration
}

func newTimeoutStep(
	name, exe, message, proj string, args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, proj, args)
	s.Timeout = timeout
	if s.Timeout <= 0 {
		s.Timeout = 30 * time.Second
	}

	return s
}

var command = exec.CommandContext

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	cmd := command(ctx, s.Exe, s.Args...)
	// Set command's working directory to the target project directory.
	cmd.Dir = s.Proj
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:  s.Name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:  s.Name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.Message, nil
}
