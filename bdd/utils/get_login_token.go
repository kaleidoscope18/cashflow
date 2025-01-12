package utils

import (
	"fmt"
	"os"
)

func GetLoginToken() (string, error) {

	password := os.Getenv("BDD_USER_PASSWORD")
	email := os.Getenv("BDD_USER_EMAIL")
	query := fmt.Sprintf(`
	query login {
		loginEmail(email: "%s", password: "%s")
	}
	`, email, password)

	var token string
	err := PostGraphQL(query, "loginEmail", "", &token)
	if err != nil {
		return "", fmt.Errorf("login failed: %v", err)
	}
	if token == "" {
		return "", fmt.Errorf("login failed: no token returned")
	}
	return token, nil
}
