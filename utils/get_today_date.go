package utils

import "time"

func GetTodayDate() string {
	now := time.Now().String()
	return ParseDate(now)
}
