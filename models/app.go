package models

type ContextKey string

const (
	JwtKey   ContextKey = "jwt"
	UserId   ContextKey = "user_id"
	UserRole ContextKey = "user_role"
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
