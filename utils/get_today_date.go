package utils

import (
	"time"
)

func GetTodayDate() string {
	location, _ := time.LoadLocation("UTC")
	now := time.Now().In(location).String()
	return ParseDate(now)
}

func GetTodayDateAsTime() time.Time {
	location, _ := time.LoadLocation("UTC")
	return time.Now().In(location)
}
