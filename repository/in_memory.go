package repository

import (
	"cashflow/models"
)

type inMemoryDatabase struct{}

var transactionById map[string]models.Transaction
var balances []models.Balance

func (database *inMemoryDatabase) init() {
	transactionById = make(map[string]models.Transaction)
	balances = make([]models.Balance, 0)
}

func (database *inMemoryDatabase) ListTransactions() []models.Transaction {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values
}

func (database *inMemoryDatabase) InsertTransaction(transaction models.Transaction) models.Transaction {
	transactionById[transaction.Id] = transaction
	return transaction
}

func (database *inMemoryDatabase) InsertBalance(amount float64, date string) models.Balance {
	newBalance := models.Balance{Date: date, Amount: amount}
	balances = append(balances, newBalance)
	return newBalance
}

func (database *inMemoryDatabase) ListBalances() []models.Balance {
	return balances
}
