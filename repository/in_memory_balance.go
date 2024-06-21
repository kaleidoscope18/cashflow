package repository

import (
	"cashflow/models"
	"context"
	"time"
)

type inMemoryBalanceDatabase struct{}

var balances []models.Balance

func (db *inMemoryBalanceDatabase) Init() error {
	balances = make([]models.Balance, 0)
	return nil
}

func (db *inMemoryBalanceDatabase) Close() error {
	balances = make([]models.Balance, 0)
	return nil
}

func (db *inMemoryBalanceDatabase) InsertBalance(amount float64, date string) (models.Balance, error) {
	newBalance := models.Balance{Date: date, Amount: amount}
	balances = append(balances, newBalance)
	return newBalance, nil
}

func (db *inMemoryBalanceDatabase) ListBalances(from time.Time, to time.Time) ([]models.Balance, error) {
	return balances, nil
}

func (db *inMemoryBalanceDatabase) DeleteBalance(ctx context.Context, date string) error {
	newList := make([]models.Balance, 0)
	for _, balance := range balances {
		if balance.Date != date {
			newList = append(newList, balance)
		}
	}

	balances = newList
	return nil
}
