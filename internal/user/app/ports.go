package app

import (
	"context"

	"blogThree/internal/user/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
	GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type PasswordPolicy interface {
	Validate(raw string) error
}

type PasswordHasher interface {
	Hash(raw string) (string, error)
	Compare(hash, raw string) error
}
