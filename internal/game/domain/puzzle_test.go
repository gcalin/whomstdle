package domain

import (
	"testing"
	"time"
)

func TestIDs(t *testing.T) {
	seen := map[PersonID]bool{}
	for _, p := range People {
		if p.ID < 1 || seen[p.ID] {
			t.Errorf("%s: ID %d is below 1 or used twice", p.Name, p.ID)
		}
		seen[p.ID] = true
	}
}

func TestSchedule(t *testing.T) {
	if Schedule[0] != 1 {
		t.Errorf("puzzle 1 has started and must stay ID 1, got %d", Schedule[0])
	}
	for i, id := range Schedule {
		if _, ok := PersonByID(id); !ok {
			t.Errorf("puzzle %d: no person with ID %d", i+1, id)
		}
		if i > 0 && Schedule[i-1] == id {
			t.Errorf("puzzles %d and %d have the same person", i, i+1)
		}
	}
	// Uses the real clock on purpose: fails a week before the schedule runs out.
	if left := len(Schedule) - int(Current(time.Now())); left < 7 {
		t.Errorf("only %d puzzles left: the schedule runs out after %s, add entries to Schedule",
			left, launch.AddDate(0, 0, len(Schedule)-1).Format("2006-01-02"))
	}
}

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

func TestPersonForIgnoresOrder(t *testing.T) {
	before := PersonFor(1).Name
	saved := People
	t.Cleanup(func() { People = saved })
	People = append([]Person{{ID: 99, Name: "New"}}, People...)
	if got := PersonFor(1).Name; got != before {
		t.Errorf("adding a person changed puzzle 1 from %s to %s", before, got)
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
