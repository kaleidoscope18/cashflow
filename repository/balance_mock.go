package repository

import (
	"cashflow/models"
)

/*
date 		| balance (transaction)
------------|----------------------
2000/01/01	| 50
2000/01/05	| 100
*/

type mockBalanceDb struct{}

func (repo *mockBalanceDb) Init()  {}
func (repo *mockBalanceDb) Close() {}

func (repo *mockBalanceDb) ListBalances() []models.Balance {
	return []models.Balance{
		{Date: "2000/01/01", Amount: 50},
		{Date: "2000/01/05", Amount: 100},
	}
}

func (repo *mockBalanceDb) InsertBalance(amount float64, date string) models.Balance {
	return models.Balance{Amount: amount, Date: date}
}
