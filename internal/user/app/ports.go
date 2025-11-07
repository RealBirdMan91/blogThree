package app

import (
	"context"

	"blogThree/internal/user/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	ExistsByEmail(ctx context.Context, email domain.Email) (bool, error)
}

type PasswordPolicy interface {
	Validate(raw string) error
}

type PasswordHasher interface {
	Hash(raw string) (string, error)
	Compare(hash, raw string) error
}
