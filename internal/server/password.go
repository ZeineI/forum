package server

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {
	var hashString string
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return hashString, fmt.Errorf("Cant generate hash: %w", err)
	}
	hashString = string(hash)
	return hashString, nil
}

func CheckPasswords(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
