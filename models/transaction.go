package models

type TransactionRepository interface {
	Init()
	Close()
	ListTransactions() []Transaction
	InsertTransaction(transaction Transaction) Transaction
}

type TransactionService interface {
	ListTransactions(todayDate *string) ([]*ComputedTransaction, error)
	WriteTransaction(date string, amount float64, description string) (*Transaction, error)
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
