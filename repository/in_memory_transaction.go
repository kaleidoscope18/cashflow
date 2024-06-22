package repository

import (
	"cashflow/models"
	"context"
	"fmt"
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

func (repo *inMemoryTransactionDatabase) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values, nil
}

func (repo *inMemoryTransactionDatabase) InsertTransaction(transaction models.Transaction) (string, error) {
	id := fmt.Sprint(len(transactionById) + 1)
	transactionById[id] = transaction
	return id, nil
}

func (repo *inMemoryTransactionDatabase) DeleteTransaction(ctx context.Context, id string) (string, error) {
	delete(transactionById, id)
	return id, nil
}
