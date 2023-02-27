package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	printTime("Start", lastAttempt)

	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()
		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			duration := time.Second * 2 << d
			shouldRetryAt := lastAttempt.Add(duration)
			printTimeWithDuration("Can retry at", shouldRetryAt, duration)

			now := time.Now()
			if !now.After(shouldRetryAt) {
				printTime("Time now", now)
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()
		response, err := circuit(ctx)
		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			printFailureAttempt(consecutiveFailures)
			consecutiveFailures++
			return response, err
		}

		printSuccessAttempt(consecutiveFailures)
		consecutiveFailures = 0
		return response, nil
	}
}

func printTime(message string, time time.Time) {
	fmt.Printf("-- %s: %s\n", message, time.Format("15:04:05.000"))
}

func printTimeWithDuration(message string, time time.Time, duration time.Duration) {
	fmt.Printf("-- %s: %s (after: %v)\n", message, time.Format("15:04:05.000"), duration)
}

func printFailureAttempt(attempt int) {
	fmt.Println(".. Attempt:", attempt, "- Failure")
}

func printSuccessAttempt(attempt int) {
	fmt.Println(".. Attempt:", attempt, "- Success")
}
