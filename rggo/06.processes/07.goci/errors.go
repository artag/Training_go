package main

import (
	"errors" // To define error values
	"fmt"    // To format messages
)

var (
	ErrValidation = errors.New("Validation failed")
	// Error when receiving a signal.
	ErrSignal = errors.New("Received signal")
)

// Error on step.
type stepErr struct {
	// Step error name.
	step string
	// Error message.
	msg string
	// Error.
	cause error
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}

	return t.step == s.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}
