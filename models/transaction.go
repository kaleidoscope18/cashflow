package models

import "time"

type TransactionService interface {
	ListTransactions(from time.Time, to time.Time) ([]ComputedTransaction, error)
	WriteTransaction(date string, amount float64, description string) (*Transaction, error)
}

type TransactionRepository interface {
	Init() error
	Close() error
	ListTransactions(from time.Time, to time.Time) ([]Transaction, error)
	InsertTransaction(transaction Transaction) (Transaction, error)
}

type Transaction struct {
	Id          string
	Amount      float64
	Date        string
	Description string
}

type Status string

const (
	StatusDone Status = "DONE"
	StatusTodo Status = "TODO"
)

type ComputedTransaction struct {
	*Transaction
	Balance float64
	Status  Status
}
