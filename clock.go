package suijin

import "time"

// Clock provides a way to get the current time.
// This interface mainly exits so that this can be mocked.
type Clock interface {
	Now() time.Time
}

// SystemClock implements the Clock interface and simply returns time.Now().
var SystemClock Clock = clock{}

type clock struct{}

func (c clock) Now() time.Time {
	return time.Now()
}
