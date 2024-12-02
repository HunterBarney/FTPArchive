package main

import (
	"fmt"
	"strings"
	"time"
)

func FormatDateTime(input string) string {
	currentTime := time.Now()
	//currentTime := time.Date(2024, 1, 5, 13, 4, 23, 0, time.UTC)
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

func toTwelveHour(hour int) int {

	if hour > 12 {
		return hour - 12
	}

	if hour == 0 {
		return 12
	}
	return hour
}

func getAMPM(hour int) string {
	if hour > 12 {
		return "PM"
	}
	return "AM"
}
