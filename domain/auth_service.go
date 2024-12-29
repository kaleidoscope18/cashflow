package domain

import (
	"cashflow/models"
	"context"
)

type authService struct {
	repository *models.Repository
}

func NewAuthService(repo *models.Repository) models.AuthService {
	s := new(authService)
	s.repository = repo
	return s
}

func (s *authService) LoginEmail(ctx context.Context, email string, password string) (string, error) {
	return (*s.repository).LoginEmail(ctx, email, password)
}

func (s *authService) ValidateSession(ctx context.Context, token string) (bool, error) {
	return (*s.repository).ValidateSession(ctx, token)
}
