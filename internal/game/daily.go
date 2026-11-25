package game

import "time"

// launch is day zero of the daily puzzle. The day changes at midnight UTC.
var launch = time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)

// TodaysPerson picks the person for the day of t, deterministically.
func TodaysPerson(t time.Time) *Person {
	days := int(t.UTC().Sub(launch).Hours() / 24)
	return &People[(days%len(People)+len(People))%len(People)]
}
