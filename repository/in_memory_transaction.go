package repository

import (
	"cashflow/models"
	"context"
	"fmt"
	"time"
)

func (repo *inMemoryRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	values := []models.Transaction{}
	for _, value := range transactionById {
		values = append(values, value)
	}

	return values, nil
}

func (repo *inMemoryRepository) InsertTransaction(transaction models.Transaction) (string, error) {
	id := fmt.Sprint(len(transactionById) + 1)
	transactionById[id] = transaction
	return id, nil
}

func (repo *inMemoryRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	delete(transactionById, id)
	return id, nil
}
