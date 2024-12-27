package repository

import (
	"log"
	"os"

	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
)

func (repo *supabaseRepository) Login(email string, password string) (string, error) {
	client := auth.New(
		os.Getenv("SUPABASE_PROJECT_REF"),
		os.Getenv("SUPABASE_API_KEY"),
	)

	session, err := client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	return session.AccessToken, nil
}

func (repo *supabaseRepository) Logout(email string) (string, error) {
	return "", nil
}
