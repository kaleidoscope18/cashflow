package models

import (
	"context"
	"time"
)

type TransactionRepository interface {
	Init() error
	Close() error
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]Transaction, error)
	InsertTransaction(transaction Transaction) (string, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
}

type TransactionService interface {
	ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]ComputedTransaction, error)
	WriteTransaction(date string, amount float64, description string, recurrency string) (string, error)
	DeleteTransaction(ctx context.Context, id string) (string, error)
	EditRecurringTransaction(ctx context.Context, editType RecurringTransactionEditType, edited TransactionEdit) (string, error)
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

type RecurringTransactionEditType string

const (
	All        RecurringTransactionEditType = "ALL"
	OnDateOnly RecurringTransactionEditType = "ON_DATE_ONLY"
	FromDate   RecurringTransactionEditType = "FROM_DATE"
)
