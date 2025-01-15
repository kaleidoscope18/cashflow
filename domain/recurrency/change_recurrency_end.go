package recurrency

import (
	"time"

	"github.com/teambition/rrule-go"
)

func ChangeRecurrencyEnd(recurrency string, newEnd time.Time) (string, error) {
	rule, err := rrule.StrToRRule(recurrency)
	if err != nil {
		return "", err
	}
	end := newEnd.AddDate(0, 0, -1)

	rule.Until(end)
	return rule.String(), nil
}
