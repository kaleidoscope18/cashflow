package models

type BalanceRepository interface {
	Init()
	Close()
	ListBalances() []Balance
	InsertBalance(amount float64, date string) Balance
}

type BalanceService interface {
	WriteBalance(balance float64, date *string) (Balance, error)
	ListBalances() ([]Balance, error)
}

// balance on a given day is the balance at the very end of this day
type Balance struct {
	Amount float64
	Date   string
}
