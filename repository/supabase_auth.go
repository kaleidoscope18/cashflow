package repository

import "log"

func (repo *supabaseRepository) LoginEmail(email string, password string) (string, error) {
	session, err := repo.db.SignInWithEmailPassword(email, password)
	if err != nil {
		log.Println("error logging in with email: ", email, err)
		return "", err
	}

	log.Println("session initialized", session)

	return session.AccessToken, nil
}

func (repo *supabaseRepository) Logout(email string) (string, error) {
	return "", nil
}
