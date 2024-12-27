package models

type App struct {
	TransactionService *TransactionService
	BalanceService     *BalanceService
	AuthService        *AuthService
}

type StorageStrategy string

const (
	Supabase StorageStrategy = "Supabase" // Supabase database instance, persisted
)

type Repository interface {
	AuthRepository
	BalanceRepository
	TransactionRepository
	Init() error
	Close() error
	Health() error
}
