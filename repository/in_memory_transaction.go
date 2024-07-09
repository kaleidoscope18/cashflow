package repository

import (
	"cashflow/models"
	"context"
	"fmt"
	"time"
)

func (repo *inMemoryRepository) GetTransactionById(ctx context.Context, id string) (models.Transaction, error) {
	return transactionById[id], nil
}

func (repo *inMemoryRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values, nil
}

func (repo *inMemoryRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	id := fmt.Sprint(len(transactionById) + 1)
	transactionById[id] = transaction
	return id, nil
}

func (repo *inMemoryRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	delete(transactionById, id)
	return id, nil
}

func (repo *inMemoryRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	transactionById[transaction.Id] = transaction
	return transaction.Id, nil
}
