package mock

import (
	"sync"
	"time"

	"github.com/Patagonicus/suijin"
)

// Clock is an implementation of suijin.Clock used for mocking.
// It starts with a given time and increases that time by a fixed amount every time Now() is called.
// Clock is safe to be used concurrently.
type Clock struct {
	current time.Time
	step    time.Duration
	lock    sync.Locker
}

// Clock has to implement suijin.Clock
var _ suijin.Clock = &Clock{}

// NewClock returns a new Clock. This Clock will return start the first time Now() is called.
// Afterwards it will return start+step, start+2*step, …
func NewClock(start time.Time, step time.Duration) *Clock {
	return &Clock{start, step, new(sync.Mutex)}
}

// Now returns a time.Time. With each call to Now the time returned is increased by the step given to NewClock.
func (c *Clock) Now() time.Time {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := c.current
	c.current = c.current.Add(c.step)

	return now
}

// Next returns the time the next call to Now will return. The difference to Now is that it will not increase the time.
func (c *Clock) Next() time.Time {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.current
}
