package models

type App struct {
	TransactionService *TransactionService
	BalanceService     *BalanceService
}

type StorageStrategy string

const (
	InMemory StorageStrategy = "InMemory" // Memory will be lost on server shutdown
	Local    StorageStrategy = "Local"    // Local mySQL database instance, persisted
	Mocked   StorageStrategy = "Mocked"   // For testing purposes, do not touch
)
