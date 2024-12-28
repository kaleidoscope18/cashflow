package domain

import (
	"cashflow/models"
)

type authService struct {
	repository *models.Repository
}

func NewAuthService(repo *models.Repository) models.AuthService {
	s := new(authService)
	s.repository = repo
	return s
}

func (s *authService) LoginEmail(email string, password string) (string, error) {
	return (*s.repository).LoginEmail(email, password)
}
