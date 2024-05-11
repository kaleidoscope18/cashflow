package utils

import (
	"strings"

	"github.com/araddon/dateparse"
)

const separator = "/"
const numberOfReplacements = 3

// ParseDate formats various date strings to YYYY/MM/DD
func ParseDate(datestring *string) string {
	t, _ := dateparse.ParseAny(*datestring)
	return strings.Replace(strings.Split(t.String(), " ")[0], "-", separator, numberOfReplacements)
}
