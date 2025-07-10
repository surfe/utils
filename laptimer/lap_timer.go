package laptimer

import (
	"sync"
	"time"
)

// Timer is an interface that defines methods for tracking the duration of multiple operations.
// It includes methods to record a lap, start the timer, and stop the timer.
type Timer interface {
	Start()
	Stop() ([]Lap, time.Duration)
	Lap(name string)
}

// Lap represents a single operation with its name and duration.
type Lap struct {
	Name     string
	Duration time.Duration
}

// LapTimer tracks the duration of multiple operations and total time.
type LapTimer struct {
	start time.Time
	laps  []Lap
	total time.Duration
	mu    sync.Mutex
}

func New() *LapTimer {
	return &LapTimer{}
}

// Lap records the duration since the last Lap call or the start.
func (t *LapTimer) Lap(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	duration := now.Sub(t.start)
	t.laps = append(t.laps, Lap{Name: name, Duration: duration})
	t.total += duration
	t.start = now // Reset start time for next lap
}

func (t *LapTimer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.start = time.Now()
}

// Stop stops the timer and returns all recorded laps and total duration.
func (t *LapTimer) Stop() ([]Lap, time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.laps, t.total
}
