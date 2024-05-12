package models

type App struct {
	TransactionService *TransactionService
	BalanceService     *BalanceService
}
