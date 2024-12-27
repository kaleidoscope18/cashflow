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

func (s *authService) Login(email string, password string) (string, error) {
	return (*s.repository).Login(email, password)
}

func (s *authService) Logout(email string) (string, error) {
	return (*s.repository).Logout(email)
}
