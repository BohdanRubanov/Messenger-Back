package utils

import (
	"errors"
	"strings"
)

func ValidateUpdateUserInput(email *string, name *string, password *string) error {
	if email != nil && strings.TrimSpace(*email) == "" {
		return errors.New("email cannot be empty")
	}
	if name != nil && strings.TrimSpace(*name) == "" {
		return errors.New("name cannot be empty")
	}
	if password != nil && strings.TrimSpace(*password) == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func ValidateCreateUserInput(email string, name string, password string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email cannot be empty")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

