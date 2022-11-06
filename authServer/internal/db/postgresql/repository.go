package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository() (*AuthRepository, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	authRep := AuthRepository{db: db}
	return &authRep, err
}

func (r *AuthRepository) checkExists(email string) (bool, error) {
	var isExists bool
	err := r.db.Get(&isExists, isUserExists, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists. Error: %s", err)
	}
	if !isExists {
		return false, fmt.Errorf("user does not exists")
	}
	return true, nil
}

func (r *AuthRepository) SignUp(email, encPassword string) error {
	var isExists bool
	err := r.db.Get(&isExists, isUserExists, email)
	if err != nil {
		return fmt.Errorf("failed to check if user exists. Error: %s", err)
	}
	if isExists {
		return fmt.Errorf("user with this email already exists")
	}
	_, err = r.db.Exec(signUp, email, encPassword)
	if err != nil {
		return fmt.Errorf("failed to sign up. Error: %s", err)
	}
	return nil
}

func (r *AuthRepository) SignIn(email, encPassword string) error {
	var isCorrectAuth bool
	err := r.db.Get(&isCorrectAuth, signIn, email, encPassword)
	if err != nil {
		return fmt.Errorf("failed to sign in. Error: %s", err)
	}
	if !isCorrectAuth {
		return fmt.Errorf("inncorrect auth data or user doesn't exist")
	}
	return nil
}
