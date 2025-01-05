package graph

import (
	"cashflow/api/graph/generated"
	"cashflow/models"
	"cashflow/utils"
	"fmt"

	"github.com/teambition/rrule-go"
)

func validateRecurrency(recurrency *string) (string, error) {
	if recurrency != nil {
		_, err := rrule.StrToRRule(*recurrency)
		if err != nil {
			return "", err
		}
		return *recurrency, err
	}
	return "", nil
}

func validateDate(date *string) *string {
	if date != nil {
		parsedDate := utils.ParseDate(*date)
		return &parsedDate
	}

	return nil
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

func validateType(input generated.TransactionEditInput) error {
	if input.Type == models.All {
		return nil
	}
	if input.Date == nil {
		return fmt.Errorf("date is required for %s type", input.Type)
	}

	return nil
}
