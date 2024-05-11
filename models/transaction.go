package models

type Status string

type TransactionService interface {
	ListTransactions(todayDate *string) []*ComputedTransaction
	WriteTransaction(date string, amount float64, description string) *Transaction
	WriteBalance(balance float64, date *string) *Balance
	ListBalances() []*Balance
}

type Transaction struct {
	Id          string
	Amount      float64
	Date        string
	Description string
}

type ComputedTransaction struct {
	*Transaction
	Balance float64
	Status  Status
}

// balance on a given day is the balance at the very end of this day
type Balance struct {
	Amount float64
	Date   string
}

const (
	StatusDone Status = "DONE"
	StatusTodo Status = "TODO"
)
