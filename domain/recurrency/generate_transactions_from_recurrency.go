package recurrency

import (
	"cashflow/models"
	"cashflow/utils"
	"fmt"
	"time"

	"github.com/teambition/rrule-go"
)

func GenerateTransactionsFromRecurrency(transactions []models.Transaction, from time.Time, to time.Time) ([]models.Transaction, error) {
	results := make([]models.Transaction, 0)
	for _, t := range transactions {
		if t.Recurrency == "" {
			return nil, fmt.Errorf("transaction with id %s does not have a recurrency", t.Id)
		}

		rule, _ := rrule.StrToRRule(t.Recurrency)
		rule.DTStart(utils.ParseDateToTime(t.Date))
		ocurrences := rule.Between(from, to, true)

		for i, occurence := range ocurrences {
			results = append(results, models.Transaction{
				Id:          fmt.Sprintf("%s-%d", t.Id, i),
				Date:        occurence.Format(utils.DateLayout),
				Amount:      t.Amount,
				Description: t.Description,
				Recurrency:  "",
			})
		}
	}

	return results, nil
}
