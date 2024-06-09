package utils

import (
	"cashflow/models"
)

func GetStatusFromDate(today string, date string) models.Status {
	if IsDateBefore(today, date) {
		return models.StatusTodo
	}
	return models.StatusDone
}
