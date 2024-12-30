package repository

import (
	"cashflow/models"
	"context"
	"errors"
	"log"
)

func (repo *supabaseRepository) LoginEmail(ctx context.Context, email string, password string) (string, error) {
	session, err := repo.db.SignInWithEmailPassword(email, password)
	if err != nil {
		log.Println("error logging in with email: ", email, err)
		return "", err
	}

	log.Println("session initialized", session)

	return session.AccessToken, nil
}

func (repo *supabaseRepository) ValidateSession(ctx context.Context, token string) (models.User, error) {
	authedClient := repo.db.Auth.WithToken(token)
	if authedClient == nil {
		return models.User{}, errors.New("invalid token")
	}
	user, err := authedClient.GetUser()

	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:    user.ID.String(),
		Role:  user.Role,
		Email: user.Email,
	}, nil
}
