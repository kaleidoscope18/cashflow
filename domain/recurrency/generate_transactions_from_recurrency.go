package recurrency

import (
	"cashflow/models"
	"cashflow/utils"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

func GenerateTransactionsFromRecurrency(transactions []models.Transaction, from time.Time, to time.Time) ([]models.Transaction, error) {
	results := make([]models.Transaction, 0)
	for _, t := range transactions {
		if t.Recurrency == "" {
			return nil, fmt.Errorf("transaction with id %s does not have a recurrency", t.Id)
		}

		var rule *rrule.RRule
		var ocurrences []time.Time
		var err error

		if len(t.Recurrency) >= 6 && t.Recurrency[0:6] == "RRULE:" {
			rules := strings.Split(t.Recurrency, "\n")
			ocurrences = make([]time.Time, 0)

			for _, ruleStr := range rules {
				rule, err = rrule.StrToRRule(ruleStr[6:]) // Skip "RRULE:" prefix
				if err != nil {
					return nil, fmt.Errorf("failed to parse recurrency rule: %w", err)
				}
				rule.DTStart(utils.ParseDateToTime(t.Date))
				ocurrences = append(ocurrences, rule.Between(from, to, true)...)
			}

			// Sort occurrences chronologically
			sort.Slice(ocurrences, func(i, j int) bool {
				return ocurrences[i].Before(ocurrences[j])
			})
		} else {
			rule, err = rrule.StrToRRule(t.Recurrency)
			if err != nil {
				return nil, fmt.Errorf("failed to parse recurrency rule: %w", err)
			}
			rule.DTStart(utils.ParseDateToTime(t.Date))
			ocurrences = rule.Between(from, to, true)
		}

		for i, occurence := range ocurrences {
			results = append(results, models.Transaction{
				Id:          fmt.Sprintf("%s-%d", t.Id, i),
				Date:        occurence.Format(utils.DateLayout),
				Amount:      t.Amount,
				Description: t.Description,
				Recurrency:  "",
				Account:     t.Account,
				IsManual:    t.IsManual,
			})
		}
	}

	return results, nil
}
