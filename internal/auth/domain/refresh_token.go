package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewRefreshToken(userID uuid.UUID, tokenHash string, ttl time.Duration) RefreshToken {
	now := time.Now().UTC()
	return RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		Revoked:   false,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}

func (t *RefreshToken) IsActive() bool {
	return !t.Revoked && !t.IsExpired()
}

func (t *RefreshToken) Revoke() {
	t.Revoked = true
}
