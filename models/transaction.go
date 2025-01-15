package models

import (
	"context"
	"time"
)

type TransactionRepository interface {
	Init() error
	Close() error
	CreateTransaction(ctx context.Context, transaction Transaction) (string, error)
	ReadTransactionById(ctx context.Context, id string, account string) (Transaction, error)
	ReadTransactions(ctx context.Context, from time.Time, to time.Time, account string) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, transaction Transaction) error
	DeleteTransaction(ctx context.Context, id string, account string) error
}

type TransactionService interface {
	ListTransactions(ctx context.Context, from time.Time, to time.Time, account string) ([]ComputedTransaction, error)
	WriteTransaction(ctx context.Context, date time.Time, amount float64, description string, recurrency string, account string) (string, error)
	DeleteTransaction(ctx context.Context, id string, account string) error
	DeleteRecurringTransaction(ctx context.Context, id string, date time.Time, account string) error
	EditTransaction(ctx context.Context, editType TransactionEditType, edited TransactionEdit) error
}

type Transaction struct {
	Id          string
	Amount      float64
	Date        string // TODO: change to time.Time
	Description string
	Recurrency  string
	Account     string
}

type TransactionEdit struct {
	Id          string
	Amount      *float64
	Date        *time.Time
	Description *string
	Recurrency  *string
	Account     string
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
