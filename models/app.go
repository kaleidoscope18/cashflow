package models

type App struct {
	TransactionService *TransactionService
	BalanceService     *BalanceService
}

type StorageStrategy string

const (
	Supabase StorageStrategy = "Supabase" // Supabase database instance, persisted
)

type Repository interface {
	BalanceRepository
	TransactionRepository
	Init() error
	Close() error
	Health() error
}
