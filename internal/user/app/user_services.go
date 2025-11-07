package app

import (
	"blogThree/internal/user/domain"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
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

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
