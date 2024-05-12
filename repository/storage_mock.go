package repository

import (
	"cashflow/models"
	"fmt"
)

/*
date 		| balance (transaction)
------------|----------------------
1999/12/31	| -10 (-10)
2000/01/01	| 50
2000/01/02	| 30 (-20)
2000/01/03	| 15 (-15)
2000/01/05	| -85 (-100)
2000/01/05	| 100
2000/01/06	| 78 (-22)
*/

type mockDb struct{}

func (repo *mockDb) Init()  {}
func (repo *mockDb) Close() {}

func (repo *mockDb) ListBalances() []models.Balance {
	return []models.Balance{
		{Date: "2000/01/01", Amount: 50},
		{Date: "2000/01/05", Amount: 100},
	}
}

func (repo *mockDb) ListTransactions() []models.Transaction {
	return []models.Transaction{
		{Date: "1999/12/31", Amount: -10},
		{Date: "2000/01/02", Amount: -20},
		{Date: "2000/01/03", Amount: -15},
		{Date: "2000/01/05", Amount: -100},
		{Date: "2000/01/06", Amount: -22},
	}
}

func (repo *mockDb) InsertBalance(amount float64, date string) models.Balance {
	return models.Balance{Amount: amount, Date: date}
}

func (repo *mockDb) InsertTransaction(transaction models.Transaction) models.Transaction {
	panic(fmt.Errorf("not implemented: InsertTransaction"))
}
