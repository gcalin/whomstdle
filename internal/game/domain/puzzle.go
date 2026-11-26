package domain

import "time"

// Puzzle is the puzzle number: 1 on launch day, then +1 each day.
type Puzzle int

// launch is the start of puzzle 1. The puzzle changes at midnight UTC.
var launch = time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)

// Current returns the puzzle number at time now. Times before launch count as puzzle 1.
func Current(now time.Time) Puzzle {
	if now.Before(launch) {
		return 1
	}
	return Puzzle(now.Sub(launch)/(24*time.Hour)) + 1
}

// Grace is how long after midnight yesterday's puzzle still accepts actions.
const Grace = time.Hour

// Start returns the time puzzle n begins.
func Start(n Puzzle) time.Time { return launch.Add(time.Duration(n-1) * 24 * time.Hour) }

// Open reports whether puzzle p accepts actions at time now: today's puzzle, or yesterday's during the grace hour.
func Open(p Puzzle, now time.Time) bool {
	c := Current(now)
	return p == c || (p == c-1 && p >= 1 && now.Before(Start(c).Add(Grace)))
}
