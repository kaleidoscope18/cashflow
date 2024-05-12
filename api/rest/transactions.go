package rest

import (
	"cashflow/models"
	"cashflow/utils"
)

type TransactionsHandler struct {
	*models.App
}

func (t TransactionsHandler) ListTransactions(today string) []*models.ComputedTransaction {
	transactions, err := t.App.TransactionService.ListTransactions(utils.GetTodayDate())
	if err != nil {
		panic(err.Error())
	}
	return t.App.TransactionService.ListTransactions(utils.GetTodayDate())
}
