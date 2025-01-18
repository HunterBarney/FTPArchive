package datetime

import (
	"fmt"
	"strings"
	"time"
)

// FormatDateTime Replaces flags an in input string with their respective values.
//
// %dd gets replaced with the current date with leading 0 if below 9.
//
// %d gets replaced with the current date without a leading 0 if below 9.
//
// %MM gets replaced with the current month with a leading 0 if below 9.
//
// %M gets replaced with the current month without a leading 0 if below 9.
//
// %yyyy gets replaced with the current year(ex. 2024).
//
// %yy gets replaced with the last 2 digits of the current year(ex. 24).
//
// %hhhh gets replaced with the current hour in 12-hour format with a leading 0 if below 9.
//
// %hhh gets replaced with the current hour in 12-hour format without a leading 0 if below 9.
//
// %hh gets replaced with the current hour in 24-hour format with a leading 0 if below 9.
//
// %h gets replaced with the current hour in 24-hour format without a leading 0 below 9.
//
// %mm gets replaced with the current minute
//
// %tt gets the current AM/PM designator
func FormatDateTime(input string, currentTime time.Time) string {
	if currentTime.IsZero() {
		currentTime = time.Now()
	}

	f := strings.NewReplacer(
		"%dd", fmt.Sprintf("%02d", currentTime.Day()),
		"%d", fmt.Sprintf("%d", currentTime.Day()),
		"%MM", fmt.Sprintf("%02d", int(currentTime.Month())),
		"%M", fmt.Sprintf("%d", int(currentTime.Month())),
		"%yyyy", fmt.Sprintf("%4d", currentTime.Year()),
		"%yy", fmt.Sprintf("%02d", currentTime.Year()%100),
		"%hhhh", fmt.Sprintf("%02d", toTwelveHour(currentTime.Hour())),
		"%hhh", fmt.Sprintf("%d", toTwelveHour(currentTime.Hour())),
		"%hh", fmt.Sprintf("%02d", currentTime.Hour()),
		"%h", fmt.Sprintf("%d", currentTime.Hour()),
		"%mm", fmt.Sprintf("%02d", currentTime.Minute()),
		"%tt", getAMPM(currentTime.Hour()),
	)

	return f.Replace(input)
}

// toTwelveHour converts the input hour to its respective 12 hour format(ex. 13 -> 1)
func toTwelveHour(hour int) int {

	if hour > 12 {
		return hour - 12
	}

	if hour == 0 {
		return 12
	}
	return hour
}

// getAMPM returns the proper AM/PM designator
func getAMPM(hour int) string {
	if hour >= 12 {
		return "PM"
	}
	return "AM"
}
