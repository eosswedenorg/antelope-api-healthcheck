package utils

import (
	"time"
)

type Time struct {
	ts time.Time
}

func (t *Time) SetTime(value time.Time) {
	t.ts = value
}

func (t Time) GetTime() time.Time {
	if !t.ts.IsZero() {
		return t.ts
	}
	return time.Now().In(time.UTC)
}
