package repository

import (
	"context"
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

func (repo *supabaseRepository) ValidateSession(ctx context.Context, token string) (bool, error) {
	authedClient := repo.db.Auth.WithToken(token)
	_, err := authedClient.GetUser()
	if err != nil {
		return false, err
	}
	return true, nil
}
