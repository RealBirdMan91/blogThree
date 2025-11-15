package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTitle   = errors.New("title is empty")
	ErrTitleTooLong = errors.New("title too long")
)

type Title struct{ s string }

func NewTitle(raw string) (Title, error) {
	t := strings.TrimSpace(raw)
	if t == "" {
		return Title{}, ErrEmptyTitle
	}
	if len([]rune(t)) > 200 {
		return Title{}, ErrTitleTooLong
	}
	return Title{s: t}, nil
}

func (t Title) String() string { return t.s }
