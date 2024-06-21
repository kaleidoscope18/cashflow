package models

import (
	"context"
	"time"
)

type TransactionRepository interface {
	Init() error
	Close() error
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]Transaction, error)
	InsertTransaction(transaction Transaction) (Transaction, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
}

type TransactionService interface {
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]ComputedTransaction, error)
	WriteTransaction(date string, amount float64, description string, recurrency string) (*Transaction, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
}

type Transaction struct {
	Id          string
	Amount      float64
	Date        string
	Description string
	Recurrency  string
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
