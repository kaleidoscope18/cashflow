package utils

import (
	"cashflow/dev"
	"time"
)

func GetTodayDate() string {
	dev.PrintJson(time.Now())
	location, _ := time.LoadLocation("UTC")
	dev.PrintJson(time.Now().In(location))
	now := time.Now().In(location).String()
	return ParseDate(now)
}
