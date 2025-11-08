package app

import (
	"blogThree/internal/auth/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, t domain.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error)
	Update(ctx context.Context, t *domain.RefreshToken) error
}

type AccessTokenEncoder interface {
	Generate(userID uuid.UUID, ttl time.Duration) (value string, expiresAt time.Time, err error)
	Verify(token string) (userID uuid.UUID, err error)
}
