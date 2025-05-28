package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func CheckPassword(existingPassword string, userPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(existingPassword), []byte(userPassword))
	return err
}
