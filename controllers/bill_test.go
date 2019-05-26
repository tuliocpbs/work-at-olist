package controllers

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestGetMonthBoundaries(t *testing.T) {
	var start time.Time
	var end time.Time
	var err error
	now := time.Now()

	_, _, err = getMonthBoundaries(int(now.Month()))
	assert.Error(t, err, "This telephone bill period has not ended")

	start, end, _ = getMonthBoundaries(int(now.Month()) - 1)
	assert.Equal(t, start.Year(), now.Year())
	assert.Equal(t, int(start.Month()), int(now.Month())-1)
	assert.Equal(t, end.Year(), now.Year())
	assert.Equal(t, int(end.Month()), int(now.Month())-1)

	start, end, _ = getMonthBoundaries(int(now.Month()) + 1)
	assert.Equal(t, start.Year(), now.Year()-1)
	assert.Equal(t, int(start.Month()), int(now.Month())+1)
	assert.Equal(t, end.Year(), now.Year()-1)
	assert.Equal(t, int(start.Month()), int(now.Month())+1)
}

func TestGetYearBoundaries(t *testing.T) {
	var start time.Time
	var end time.Time
	var err error
	now := time.Now()

	_, _, err = getYearBoundaries(now.Year())
	assert.Error(t, err, "This telephone bill period has not ended")

	_, _, _ = getYearBoundaries(now.Year() + 1)
	assert.Error(t, err, "This telephone bill period has not ended")

	start, end, _ = getYearBoundaries(now.Year() - 1)
	firstDay := time.Date(now.Year()-1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(1, 0, 0).Add(-time.Nanosecond)
	assert.Equal(t, start, firstDay)
	assert.Equal(t, end, lastDay)

}
