package rest

import (
	"cashflow/models"
)

type TransactionsHandler struct {
	*models.App
}

func (t TransactionsHandler) ListTransactions(today string) []*models.ComputedTransaction {
	transactions, err := (*t.App.TransactionService).ListTransactions(nil)
	if err != nil {
		panic(err.Error())
	}
	return transactions
}
