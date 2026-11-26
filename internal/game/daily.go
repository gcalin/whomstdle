package game

import "time"

// launch is the start of puzzle 1. The puzzle changes at midnight UTC.
var launch = time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)

// Schedule[n-1] is the person ID for puzzle n.
// Only add to the end. Never change an entry for a puzzle that has started.
// Do not schedule people in ID order, or the IDs in the page give away the answer.
var Schedule = []int{1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2, 1, 3, 2}

// Current returns the puzzle number at time now. Times before launch count as puzzle 1.
func Current(now time.Time) int {
	if now.Before(launch) {
		return 1
	}
	return int(now.Sub(launch)/(24*time.Hour)) + 1
}

// PersonFor returns the person for puzzle n.
func PersonFor(n int) *Person {
	id := Schedule[(n-1)%len(Schedule)] // wraps if the schedule runs out; a test warns a week before
	p, _ := PersonByID(id)
	return p
}
