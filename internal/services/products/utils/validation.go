package utils

import (
	"errors"
	"strings"
)

func ValidateUpdateProductInput(title *string, description *string, price *int) error {
	if title != nil && strings.TrimSpace(*title) == "" {
		return errors.New("title cannot be empty")
	}
	if description != nil && strings.TrimSpace(*description) == "" {
		return errors.New("description cannot be empty")
	}
	if price != nil && *price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
}

func ValidateCreateProductInput(title string, description string, price int) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	}
	if strings.TrimSpace(description) == "" {
		return errors.New("description cannot be empty")
	}
	if price <= 0 {
		return errors.New("price must be greater than zero")
	}
	return nil
}
