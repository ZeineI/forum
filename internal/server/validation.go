package server

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

func CheckEmail(email string) (bool, error) {
	isCorrect, err := regexp.MatchString(`^([\w\.\_]{2,30})@(\w{1,}).([a-z]{2,4})$`, email)
	if err != nil {
		log.Println("invalid email")
		return isCorrect, fmt.Errorf("invalid email validation: %w", err)
	}
	if !isCorrect {
		log.Println("invalid email - bool")
		return false, nil
	}
	return true, nil
}

func WitchNotUnique(err string) string {
	if strings.Contains(err, "User.email") {
		return "Email already exist"
	} else {
		return "Username already exist"
	}
}

func IsValidPassword(password string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
