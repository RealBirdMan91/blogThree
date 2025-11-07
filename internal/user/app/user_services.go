package app

import (
	"blogThree/internal/user/domain"
	"context"
	"errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email already in use")
)

type Service struct {
	repo   UserRepository
	policy PasswordPolicy
	hasher PasswordHasher
}

type UserService interface {
	SignUp(ctx context.Context, email, rawPassword string) (*domain.User, error)
}

func NewService(r UserRepository, policy PasswordPolicy, h PasswordHasher) UserService {
	return &Service{
		repo:   r,
		policy: policy,
		hasher: h,
	}
}

func (s *Service) SignUp(ctx context.Context, email, rawPassword string) (*domain.User, error) {
	emailVO, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByEmail(ctx, emailVO)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	if err := s.policy.Validate(rawPassword); err != nil {
		return nil, err
	}

	hashedPassword, err := s.hasher.Hash(rawPassword)
	if err != nil {
		return nil, err
	}
	passwordHashVO, err := domain.NewPasswordHash(hashedPassword)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(emailVO, passwordHashVO)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
