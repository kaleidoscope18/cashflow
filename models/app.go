package models

type ContextKey string

const (
	JwtKey ContextKey = "jwt"
)

type App struct {
	TransactionService *TransactionService
	BalanceService     *BalanceService
	AuthService        *AuthService
}

type StorageStrategy string

const (
	Supabase StorageStrategy = "Supabase"
)

type Repository interface {
	AuthRepository
	BalanceRepository
	TransactionRepository
	Init() error
	Close() error
}
