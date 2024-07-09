package models

import (
	"context"
	"time"
)

type TransactionRepository interface {
	Init() error
	Close() error
	GetTransactionById(ctx context.Context, id string) (Transaction, error)
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]Transaction, error)
	InsertTransaction(ctx context.Context, transaction Transaction) (string, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
	EditTransaction(ctx context.Context, transaction Transaction) (string, error)
}

type TransactionService interface {
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]ComputedTransaction, error)
	WriteTransaction(ctx context.Context, date string, amount float64, description string, recurrency string) (string, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
	EditTransaction(ctx context.Context, editType TransactionEditType, edited TransactionEdit) (string, error)
}

type Transaction struct {
	Id          string
	Amount      float64
	Date        string
	Description string
	Recurrency  string
}

type TransactionEdit struct {
	Id          string
	Amount      *float64
	Date        *string
	Description *string
	Recurrency  *string
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

type TransactionEditType string

const (
	All        TransactionEditType = "ALL"
	OnDateOnly TransactionEditType = "ON_DATE_ONLY"
	FromDate   TransactionEditType = "FROM_DATE"
)
