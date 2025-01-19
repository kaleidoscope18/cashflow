package graph

import (
	"cashflow/api/graph/generated"
	"cashflow/models"
	"cashflow/utils"
	"fmt"

	"github.com/teambition/rrule-go"
)

func validateRecurrency(recurrency *string) error {
	if recurrency != nil {
		_, err := rrule.StrToRRule(*recurrency)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func validateDate(date *string, required bool) error {
	if date != nil {
		parsedDate := utils.ParseDate(*date)
		if parsedDate == "0001-01-01" {
			return fmt.Errorf("invalid date: %s", *date)
		}
	}

	if required && date == nil {
		return fmt.Errorf("date is required")
	}

	return nil
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

func validatePlanRecurrency(recurrency *string) error {
	if recurrency != nil {
		rule, err := rrule.StrToRRule(*recurrency)
		if err != nil {
			return err
		}
		if rule.OrigOptions.Count == 0 {
			return fmt.Errorf("should have a minimum of 1 occurrence: %s", *recurrency)
		}
	}
	return nil
}
