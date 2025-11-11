package app

import (
	apperr "blogThree/internal/errors"
	"blogThree/internal/user/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo   UserRepository
	policy PasswordPolicy
	hasher PasswordHasher
}

type UserService interface {
	SignUp(ctx context.Context, email, rawPassword string) (*domain.User, error)
	SignIn(ctx context.Context, email, rawPassword string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
}

func NewService(r UserRepository, policy PasswordPolicy, h PasswordHasher) UserService {
	return &Service{
		repo:   r,
		policy: policy,
		hasher: h,
	}
}

func (s *Service) SignIn(ctx context.Context, email, rawPassword string) (*domain.User, error) {
	emailVO, err := domain.NewEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.GetByEmail(ctx, emailVO)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.hasher.Compare(user.PasswordHash().String(), rawPassword); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *Service) SignUp(ctx context.Context, email, rawPassword string) (*domain.User, error) {
	emailVO, err := domain.NewEmail(email)
	if err != nil {
		return nil, NewInvalidEmailError()
	}

	exists, err := s.repo.ExistsByEmail(ctx, emailVO)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	if err := s.policy.Validate(rawPassword); err != nil {
		return nil, NewWeakPasswordError(err.Error())
	}

	hashedPassword, err := s.hasher.Hash(rawPassword)
	if err != nil {
		return nil, NewPasswordHashFailed(err)
	}
	passwordHashVO, err := domain.NewPasswordHash(hashedPassword)
	if err != nil {
		return nil, NewPasswordHashFailed(err)
	}

	user, err := domain.NewUser(emailVO, passwordHashVO)
	if err != nil {
		return nil, apperr.Unknown(err)
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}
