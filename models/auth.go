package models

import "context"

type AuthService interface {
	LoginEmail(ctx context.Context, email string, password string) (string, error)
	ValidateSession(ctx context.Context, token string) (bool, error)
}

type AuthRepository interface {
	LoginEmail(ctx context.Context, email string, password string) (string, error)
	ValidateSession(ctx context.Context, token string) (bool, error)
}
