package recurrency

import "cashflow/models"

func SplitTransactionsWithRecurrency(transactions []models.Transaction) ([]models.Transaction, []models.Transaction) {
	var withoutRecurrency []models.Transaction
	var withRecurrency []models.Transaction

	for _, transaction := range transactions {
		if transaction.Recurrency == "" {
			withoutRecurrency = append(withoutRecurrency, transaction)
		} else {
			withRecurrency = append(withRecurrency, transaction)
		}
	}

	return withRecurrency, withoutRecurrency
}
