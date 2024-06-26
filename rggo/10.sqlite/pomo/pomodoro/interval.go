package pomodoro

import (
	"context" // To carry context and cancellation signals from the user interface
	"errors"  // To define custom errors
	"fmt"     // To format output
	"time"    // To handle time-related data
)

// Category constants
const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

// State constants
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Pomodoro interval
type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type Repository interface {
	// Create an interval
	Create(i Interval) (int64, error)
	// Update the interval
	Update(i Interval) error
	// Retrieve an interval by id
	ByID(id int64) (Interval, error)
	// Find the last interval
	Last() (Interval, error)
	// Retrieve intervals of type break
	Breaks(n int) ([]Interval, error)
	// Return a daily summary
	CategorySummary(day time.Time, filter string) (time.Duration, error)
}

// Errors
var (
	//lint:ignore ST1005 Ignore warning
	ErrNoIntervals        = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted  = errors.New("Interval is completed or cancelled")
	//lint:ignore ST1005 Ignore warning
	ErrInvalidState = errors.New("Invalid State")
	//lint:ignore ST1005 Ignore warning
	ErrInvalidID = errors.New("Invalid ID")
)

// Configuration required to instatiate an interval
type IntervalConfig struct {
	repo               Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

// Instatiate a new IntervalConfig
func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {
	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}

	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}
	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}
	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}

	return c
}

// Returns the next interval category
func nextCategory(r Repository) (string, error) {
	li, err := r.Last()
	if err != nil && err == ErrNoIntervals {
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}

	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}

	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}

	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}

	return CategoryLongBreak, nil
}

// Callback function
type Callback func(Interval)

// Control the interval timer
func tick(
	ctx context.Context,
	id int64,
	config *IntervalConfig,
	start, periodic, end Callback) error {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	i, err := config.repo.ByID(id)
	if err != nil {
		return err
	}

	expire := time.After(i.PlannedDuration - i.ActualDuration)
	start(i)

	for {
		select {
		case <-ticker.C:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}

			if i.State == StatePaused {
				return nil
			}

			i.ActualDuration += time.Second
			if err := config.repo.Update(i); err != nil {
				return err
			}

			periodic(i)

		case <-expire:
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}

			i.State = StateDone
			end(i)
			return config.repo.Update(i)

		case <-ctx.Done():
			i, err := config.repo.ByID(id)
			if err != nil {
				return err
			}

			i.State = StateCancelled
			return config.repo.Update(i)
		}
	}
}

// Create new Interval instance
func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return i, err
	}

	i.Category = category

	switch category {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}

	if i.ID, err = config.repo.Create(i); err != nil {
		return i, err
	}

	return i, nil
}

// Retrieve the last interval from the repository
func GetInterval(config *IntervalConfig) (Interval, error) {
	var i Interval
	var err error

	i, err = config.repo.Last()
	if err != nil && err != ErrNoIntervals {
		return i, err
	}
	if err == nil && i.State != StateCancelled && i.State != StateDone {
		return i, nil
	}

	return newInterval(config)
}

// Start the interval timer
func (i Interval) Start(
	ctx context.Context,
	config *IntervalConfig,
	start, periodic, end Callback) error {

	switch i.State {
	case StateRunning:
		return nil
	case StateNotStarted:
		i.StartTime = time.Now()
		fallthrough
	case StatePaused:
		i.State = StateRunning
		if err := config.repo.Update(i); err != nil {
			return err
		}
		return tick(ctx, i.ID, config, start, periodic, end)
	case StateCancelled, StateDone:
		return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
	default:
		return fmt.Errorf("%w: %d", ErrInvalidState, i.State)
	}
}

// Verifies whether the instance of Interval is running and pauses it
func (i Interval) Pause(config *IntervalConfig) error {
	if i.State != StateRunning {
		return ErrIntervalNotRunning
	}

	i.State = StatePaused
	return config.repo.Update(i)
}
