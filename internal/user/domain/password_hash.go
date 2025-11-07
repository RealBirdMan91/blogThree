package domain

import (
	"errors"
	"strings"
)

var ErrEmptyPasswordHash = errors.New("password hash is empty")

type PasswordHash struct{ s string }

func NewPasswordHash(raw string) (PasswordHash, error) {
	if strings.TrimSpace(raw) == "" {
		return PasswordHash{}, ErrEmptyPasswordHash
	}
	return PasswordHash{s: raw}, nil
}

func (p PasswordHash) String() string { return p.s }
