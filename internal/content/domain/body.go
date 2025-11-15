package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyBody = errors.New("body is empty")
)

type Body struct{ s string }

func NewBody(raw string) (Body, error) {
	b := strings.TrimSpace(raw)
	if b == "" {
		return Body{}, ErrEmptyBody
	}
	return Body{s: b}, nil
}

func (b Body) String() string { return b.s }
