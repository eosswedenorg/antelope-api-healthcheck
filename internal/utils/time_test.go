package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_GetTimeWithDefaultValue(t *testing.T) {
	var ts Time

	// Assert that time is NOW (+-10 seconds)
	assert.InDelta(t, ts.GetTime().Unix(), time.Now().In(time.UTC).Unix(), float64(10))
}

func TestTime_GetTimeWithSetTime(t *testing.T) {
	var ts Time

	expected := time.Unix(1048722042, 500)
	ts.SetTime(expected)

	assert.Equal(t, expected, ts.GetTime())
}
