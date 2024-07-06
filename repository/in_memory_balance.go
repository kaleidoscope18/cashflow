package repository

import (
	"cashflow/models"
	"context"
	"time"
)

func (db *inMemoryRepository) InsertBalance(amount float64, date string) (models.Balance, error) {
	newBalance := models.Balance{Date: date, Amount: amount}
	balances = append(balances, newBalance)
	return newBalance, nil
}

func (db *inMemoryRepository) ListBalances(from time.Time, to time.Time) ([]models.Balance, error) {
	return balances, nil
}

func (db *inMemoryRepository) DeleteBalance(ctx context.Context, date string) error {
	newList := make([]models.Balance, 0)
	for _, balance := range balances {
		if balance.Date != date {
			newList = append(newList, balance)
		}
	}

	balances = newList
	return nil
}
