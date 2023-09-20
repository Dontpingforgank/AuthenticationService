package PasswordUtils

import (
	"github.com/Dontpingforgank/AuthenticationService/CustomErrors"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func VerifyPassword(password string) (bool, error) {
	if len(password) >= 6 {
		var upper bool
		var lower bool
		var number bool

		for _, char := range password {
			if unicode.IsLower(char) {
				lower = true
				continue
			}

			if unicode.IsUpper(char) {
				upper = true
				continue
			}

			if unicode.IsNumber(char) {
				number = true
			}
		}

		if upper && lower && number {
			return true, nil
		} else {
			return false, CustomErrors.InsecurePassword{}
		}
	}

	return false, nil
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}
