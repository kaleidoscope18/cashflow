package repository

import (
	"cashflow/models"
)

type inMemoryBalanceDatabase struct{}

var balances []models.Balance

func (db *inMemoryBalanceDatabase) Init() {
	balances = make([]models.Balance, 0)
}

func (db *inMemoryBalanceDatabase) Close() {}

func (db *inMemoryBalanceDatabase) InsertBalance(amount float64, date string) models.Balance {
	newBalance := models.Balance{Date: date, Amount: amount}
	balances = append(balances, newBalance)
	return newBalance
}

func (db *inMemoryBalanceDatabase) ListBalances() []models.Balance {
	return balances
}
