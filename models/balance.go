package models

import "time"

type BalanceRepository interface {
	Init() error
	Close() error
	ListBalances(from time.Time, to time.Time) ([]Balance, error)
	InsertBalance(amount float64, date string) (Balance, error)
}

type BalanceService interface {
	WriteBalance(balance float64, date *string) (Balance, error)
	ListBalances(from time.Time, to time.Time) ([]Balance, error)
}

// balance on a given day is the balance at the very end of this day
type Balance struct {
	Amount float64
	Date   string
}
