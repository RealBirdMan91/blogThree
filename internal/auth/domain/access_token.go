package domain

import "time"

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

func NewAccessToken(value string, expiresAt time.Time) AccessToken {
	return AccessToken{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func (t AccessToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
