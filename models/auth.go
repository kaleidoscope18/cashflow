package models

type AuthRepository interface {
	Login(email string, password string) (string, error)
	Logout(email string) (string, error)
}

type AuthService interface {
	Login(email string, password string) (string, error)
	Logout(email string) (string, error)
}
