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
	PlanService        *PlanService
}

type StorageStrategy string

const (
	Supabase StorageStrategy = "Supabase"
)

type Repository interface {
	AuthRepository
	BalanceRepository
	TransactionRepository
	PlanRepository
	Init() error
	Close() error
}
