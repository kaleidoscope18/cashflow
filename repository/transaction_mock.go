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

type mockTransactionDb struct{}

func (repo *mockTransactionDb) Init()  {}
func (repo *mockTransactionDb) Close() {}

func (repo *mockTransactionDb) ListTransactions() []models.Transaction {
	return []models.Transaction{
		{Date: "1999/12/31", Amount: -10},
		{Date: "2000/01/02", Amount: -20},
		{Date: "2000/01/03", Amount: -15},
		{Date: "2000/01/05", Amount: -100},
		{Date: "2000/01/06", Amount: -22},
	}
}

func (repo *mockTransactionDb) InsertTransaction(transaction models.Transaction) models.Transaction {
	panic(fmt.Errorf("not implemented: InsertTransaction"))
}
