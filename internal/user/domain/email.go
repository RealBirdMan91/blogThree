package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyEmail      = errors.New("email is empty")
	ErrInvalidEmailFmt = errors.New("invalid email format")
)

type Email struct{ s string }

func NewEmail(raw string) (Email, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Email{}, ErrEmptyEmail
	}
	if !isValidEmail(trimmed) {
		return Email{}, ErrInvalidEmailFmt
	}
	return Email{s: strings.ToLower(trimmed)}, nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (email Email) String() string {
	return email.s
}
