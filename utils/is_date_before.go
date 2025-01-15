package utils

import "time"

const DateLayout = "2006-01-02" // (year, month, day)

func IsDateBefore(earlierDate string, laterDate string) bool {
	date1 := ParseDateToTime(earlierDate)
	date2 := ParseDateToTime(laterDate)

	return date1.Before(date2)
}

func IsTimeBeforeTime(earlierDate time.Time, laterDate time.Time) bool {
	return earlierDate.Before(laterDate)
}
