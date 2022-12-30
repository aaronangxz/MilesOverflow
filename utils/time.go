package utils

import "time"

const (
	//TimeZone constant timezone
	TimeZone = "Asia/Singapore"
)

var (
	tz, _ = time.LoadLocation(TimeZone)
)

// UnixToUTC Converts current unix time to UTC time object
func UnixToUTC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Local().UTC()
}

// MonthStartEndDate Returns the start and end day of the current month in SGT unix time
func MonthStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)
	currentYear, currentMonth, _ := date.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, tz)
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 23, 59, 59, 59, tz)
	return firstOfMonth.Unix(), lastOfMonth.Unix()
}
