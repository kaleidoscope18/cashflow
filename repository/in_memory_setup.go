package repository

import "cashflow/models"

type inMemoryRepository struct {
}

var balances []models.Balance
var transactionById map[string]models.Transaction

func (db *inMemoryRepository) Init() error {
	balances = make([]models.Balance, 0)
	transactionById = make(map[string]models.Transaction)
	return nil
}

func (db *inMemoryRepository) Close() error {
	balances = make([]models.Balance, 0)
	transactionById = make(map[string]models.Transaction)
	return nil
}

func (db *inMemoryRepository) Health() error {
	return nil
}
