package stats

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

// A Timer that can be started and stopped and accumulates the total time it
// was running (the time between Start() and Stop()).
type Timer struct {
	name     fmt.Stringer
	created  int
	start    time.Time
	duration time.Duration
}

// Start the timer.
func (t *Timer) Start() *Timer {
	t.start = time.Now()
	return t
}

// Stop the timer.
func (t *Timer) Stop() {
	t.duration += time.Since(t.start)
}

// ElapsedTime returns the time that passed since starting the timer.
func (t *Timer) ElapsedTime() time.Duration {
	return time.Since(t.start)
}

// Duration returns the duration value of the timer in seconds.
func (t *Timer) Duration() float64 {
	return t.duration.Seconds()
}

// Return a string representation of the Timer.
func (t *Timer) String() string {
	return fmt.Sprintf("%s: %s", t.name, t.duration)
}

// A TimerGroup represents a group of timers relevant to a single query.
type TimerGroup struct {
	timers map[fmt.Stringer]*Timer
}

// NewTimerGroup constructs a new TimerGroup.
func NewTimerGroup() *TimerGroup {
	return &TimerGroup{timers: map[fmt.Stringer]*Timer{}}
}

// GetTimer gets (and creates, if necessary) the Timer for a given code section.
func (t *TimerGroup) GetTimer(name fmt.Stringer) *Timer {
	if timer, exists := t.timers[name]; exists {
		return timer
	}
	timer := &Timer{
		name:    name,
		created: len(t.timers),
	}
	t.timers[name] = timer
	return timer
}

// Timers is a slice of Timer pointers that implements Len and Swap from
// sort.Interface.
type Timers []*Timer

type byCreationTimeSorter struct{ Timers }

// Len implements sort.Interface.
func (t Timers) Len() int {
	return len(t)
}

// Swap implements sort.Interface.
func (t Timers) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (s byCreationTimeSorter) Less(i, j int) bool {
	return s.Timers[i].created < s.Timers[j].created
}

// Return a string representation of a TimerGroup.
func (t *TimerGroup) String() string {
	timers := byCreationTimeSorter{}
	for _, timer := range t.timers {
		timers.Timers = append(timers.Timers, timer)
	}
	sort.Sort(timers)
	result := &bytes.Buffer{}
	for _, timer := range timers.Timers {
		fmt.Fprintf(result, "%s\n", timer)
	}
	return result.String()
}
