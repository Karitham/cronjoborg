package cronjoborg

import (
	"time"
)

// Seconds represents a time.Time in seconds
type Seconds int64

// Time returns a time.Time from the Seconds
func (s Seconds) Time() time.Time {
	return time.Unix(int64(s), 0)
}

// Microseconds represents a time.Time in microseconds
type Microseconds int64

// Time returns a time.Time from the Microseconds
func (m Microseconds) Time() time.Time {
	return time.UnixMicro(int64(m))
}

// Milliseconds represents a time.Time in milliseconds
type Milliseconds int64

func (m Milliseconds) Time() time.Time {
	return time.UnixMilli(int64(m))
}
