package recurrency

import (
	"cashflow/models"
	"log"

	"github.com/teambition/rrule-go"
)

func GenerateTransactionFromPlan(planName string, account string, startDate string, recurrency string, totalAmount float64) (models.Transaction, error) {
	rule, err := rrule.StrToRRule(recurrency)
	if err != nil {
		log.Printf("error generating transaction from plan: %v", err)
		return models.Transaction{}, err
	}
	numberOfTransactions := float64((*rule).Options.Count)

	transaction := models.Transaction{
		Amount:      float64(int((totalAmount/numberOfTransactions)*100)) / 100, // Rounded to nearest 0.01
		Date:        startDate,
		Description: "from plan: " + planName,
		Recurrency:  recurrency,
		Account:     account,
	}

	return transaction, nil
}
