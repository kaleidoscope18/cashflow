package validation

import (
	"cashflow/utils"
	"fmt"
)

func VerifyDate(date *string) (*string, error) {
	if date != nil && *date != "" {
		if !utils.IsValidDate(*date) {
			return nil, fmt.Errorf("invalid date")
		}

		parsedDate := utils.ParseDate(*date)

		return &parsedDate, nil
	}

	return nil, nil
}
