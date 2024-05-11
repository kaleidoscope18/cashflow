package utils

import (
	"time"
)

const Layout = "2006/01/02" // (year, month, day)

func IsDateBefore(earlierDate string, laterDate string) bool {
	date1, _ := time.Parse(Layout, earlierDate)
	date2, _ := time.Parse(Layout, laterDate)

	return date1.Before(date2)
}
