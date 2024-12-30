package models

import "context"

type User struct {
	ID    string
	Role  string
	Email string
}

type AuthService interface {
	LoginEmail(ctx context.Context, email string, password string) (string, error)
	ValidateSession(ctx context.Context, token string) (User, error)
}

type AuthRepository interface {
	LoginEmail(ctx context.Context, email string, password string) (string, error)
	ValidateSession(ctx context.Context, token string) (User, error)
}
