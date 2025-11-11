package policies

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrTooShort = errors.New("password too short")
	ErrNoLetter = errors.New("password must contain at least one letter")
)

type SimplePasswordPolicy struct {
	MinLen        int
	RequireLetter bool
}

func NewSimplePasswordPolicy(minLen int, requireLetter bool) SimplePasswordPolicy {
	return SimplePasswordPolicy{MinLen: minLen, RequireLetter: requireLetter}
}

var letterRegexp = regexp.MustCompile(`[A-Za-z]`)

func (p SimplePasswordPolicy) Validate(raw string) error {
	s := strings.TrimSpace(raw)

	if len(s) < p.MinLen {
		return ErrTooShort
	}
	if p.RequireLetter && !letterRegexp.MatchString(s) {
		return ErrNoLetter
	}
	return nil
}
