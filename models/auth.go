package models

type AuthRepository interface {
	LoginEmail(email string, password string) (string, error)
}

type AuthService interface {
	LoginEmail(email string, password string) (string, error)
}
