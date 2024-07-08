package utils

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

const separator = "/"
const numberOfReplacements = 3

// ParseDate formats various date strings to YYYY/MM/DD
func ParseDate(datestring string) string {
	t, _ := dateparse.ParseAny(datestring)
	return strings.Replace(strings.Split(t.UTC().String(), " ")[0], "-", separator, numberOfReplacements)
}

func ParseDateToTime(datestring string) time.Time {
	t, _ := dateparse.ParseAny(datestring)
	return t.UTC()
}
