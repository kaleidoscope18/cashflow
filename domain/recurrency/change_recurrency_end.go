package recurrency

import (
	"cashflow/utils"

	"github.com/teambition/rrule-go"
)

func ChangeRecurrencyEnd(recurrency string, newEnd string) (string, error) {
	rule, err := rrule.StrToRRule(recurrency)
	if err != nil {
		return "", err
	}
	end := utils.ParseDateToTime(newEnd).AddDate(0, 0, -1)

	rule.Until(end)
	return rule.String(), nil
}
