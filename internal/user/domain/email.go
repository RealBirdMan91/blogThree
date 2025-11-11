package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyEmail   = errors.New("email is empty")
	ErrInvalidEmail = errors.New("invalid email")
)

type Email struct{ s string }

func NewEmail(raw string) (Email, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Email{}, ErrEmptyEmail
	}
	if !isValidEmail(trimmed) {
		return Email{}, ErrInvalidEmail
	}
	return Email{s: strings.ToLower(trimmed)}, nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool { return emailRegex.MatchString(email) }

func (e Email) String() string { return e.s }
