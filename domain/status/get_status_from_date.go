package status

import (
	"cashflow/models"
	"cashflow/utils"
)

func GetStatusFromDate(today string, date string) models.Status {
	if utils.IsDateBefore(today, date) {
		return models.StatusTodo
	}
	return models.StatusDone
}
