package policies

import (
	"blogThree/internal/user/app"
	"errors"
	"regexp"
	"strings"
)

type SimplePasswordPolicy struct {
	MinLen        int
	RequireLetter bool
}

var _ app.PasswordPolicy = (*SimplePasswordPolicy)(nil)

func NewSimplePasswordPolicy(minLen int, requireLetter bool) SimplePasswordPolicy {
	return SimplePasswordPolicy{MinLen: minLen, RequireLetter: requireLetter}
}

var letterRegexp = regexp.MustCompile(`[A-Za-z]`)

func (p SimplePasswordPolicy) Validate(raw string) error {
	s := strings.TrimSpace(raw)
	if len(s) < p.MinLen {
		return errors.New("password too short")
	}
	if p.RequireLetter && !letterRegexp.MatchString(s) {
		return errors.New("password must contain at least one letter")
	}
	return nil
}
