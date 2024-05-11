package models

type Repository interface {
	ListTransactions() []Transaction
	InsertTransaction(transaction Transaction) Transaction

	ListBalances() []Balance
	InsertBalance(amount float64, date string) Balance
}
