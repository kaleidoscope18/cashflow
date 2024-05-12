package repository

import (
	"cashflow/models"
)

type inMemoryTransactionDatabase struct{}

var transactionById map[string]models.Transaction

func (repo *inMemoryTransactionDatabase) Init() {
	transactionById = make(map[string]models.Transaction)
}

func (repo *inMemoryTransactionDatabase) Close() {}

func (repo *inMemoryTransactionDatabase) ListTransactions() []models.Transaction {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values
}

func (repo *inMemoryTransactionDatabase) InsertTransaction(transaction models.Transaction) models.Transaction {
	transactionById[transaction.Id] = transaction
	return transaction
}
