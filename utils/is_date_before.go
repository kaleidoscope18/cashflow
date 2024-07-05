package utils

import (
	"time"
)

const DateLayout = "2006/01/02" // (year, month, day)

func IsDateBefore(earlierDate string, laterDate string) bool {
	date1, _ := time.Parse(DateLayout, earlierDate)
	date2, _ := time.Parse(DateLayout, laterDate)

	return date1.Before(date2)
}
