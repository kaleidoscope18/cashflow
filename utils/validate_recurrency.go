package utils

import (
	"github.com/teambition/rrule-go"
)

func ValidateRecurrency(recurrency string) error {
	_, err := rrule.StrToRRule(recurrency)
	return err
}
