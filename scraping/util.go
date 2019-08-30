package scraping

import "time"

var location, _ = time.LoadLocation("America/New_York")

func processDayDifference(baseTime *time.Time, hour, min int) time.Time {
	fixedTime := time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), hour, min, 0, 0, location)
	if hour < baseTime.Hour() {
		fixedTime = time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day()+1, hour, min, 0, 0, location)
	}
	return fixedTime
}
