package domain

import (
	"testing"
	"time"
)

func TestCurrent(t *testing.T) {
	for _, c := range []struct {
		at   time.Time
		want Puzzle
	}{
		{launch.Add(-time.Minute), 1},
		{launch, 1},
		{launch.Add(24*time.Hour - time.Second), 1},
		{launch.Add(24 * time.Hour), 2},
	} {
		if got := Current(c.at); got != c.want {
			t.Errorf("Current(%v) = %d, want %d", c.at, got, c.want)
		}
	}
}

func TestOpen(t *testing.T) {
	day2 := Start(2)
	for _, c := range []struct {
		p    Puzzle
		at   time.Time
		want bool
	}{
		{2, day2, true},                       // today
		{1, day2.Add(59 * time.Minute), true}, // yesterday, in the grace hour
		{1, day2.Add(61 * time.Minute), false},
		{3, day2, false}, // tomorrow
		{0, launch, false},
	} {
		if got := Open(c.p, c.at); got != c.want {
			t.Errorf("Open(%d, %v) = %v, want %v", c.p, c.at, got, c.want)
		}
	}
}
