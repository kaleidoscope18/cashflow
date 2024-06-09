package repository

import (
	"cashflow/models"
	"time"
)

type inMemoryTransactionDatabase struct{}

var transactionById map[string]models.Transaction

func (repo *inMemoryTransactionDatabase) Init() error {
	transactionById = make(map[string]models.Transaction)
	return nil
}

func (repo *inMemoryTransactionDatabase) Close() error {
	return nil
}

func (repo *inMemoryTransactionDatabase) ListTransactions(from time.Time, to time.Time) ([]models.Transaction, error) {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values, nil
}

func (repo *inMemoryTransactionDatabase) InsertTransaction(transaction models.Transaction) (models.Transaction, error) {
	transactionById[transaction.Id] = transaction
	return transaction, nil
}
