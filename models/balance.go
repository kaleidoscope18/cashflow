package models

import (
	"context"
	"time"
)

type BalanceRepository interface {
	InsertBalance(ctx context.Context, amount float64, date string) (Balance, error)
	ListBalances(ctx context.Context, from time.Time, to time.Time) ([]Balance, error)
	DeleteBalance(ctx context.Context, date string) error
}

type BalanceService interface {
	WriteBalance(ctx context.Context, balance float64, date *string) (Balance, error)
	ListBalances(ctx context.Context, from time.Time, to time.Time) ([]Balance, error)
	DeleteBalance(ctx context.Context, date string) error
}

// balance on a given day is the balance at the very end of this day
type Balance struct {
	Amount float64
	Date   string
}
