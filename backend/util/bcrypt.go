package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordEmpty = fmt.Errorf("password is empty")

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed hash password")
	}
	return string(bytes), nil
}
