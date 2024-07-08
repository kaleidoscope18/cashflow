package graph

import (
	"cashflow/utils"

	"github.com/teambition/rrule-go"
)

func validateRecurrency(recurrency *string) (string, error) {
	if recurrency != nil {
		_, err := rrule.StrToRRule(*recurrency)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func validateDate(date *string) *string {
	if date != nil {
		utils.ParseDate(*date)
	}
	return date
}

func validateDescription(description *string) string {
	var d string
	if description == nil {
		d = ""
	} else {
		d = *description
	}

	return d
}
