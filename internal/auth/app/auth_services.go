package app

import (
	"blogThree/internal/auth/domain"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  domain.AccessToken
	RefreshToken string // raw, kommt ins HttpOnly-Cookie
}

type AuthService interface {
	GenerateForUser(ctx context.Context, userID uuid.UUID) (*Tokens, error)
	Refresh(ctx context.Context, rawRefresh string) (*Tokens, uuid.UUID, error)
	Revoke(ctx context.Context, rawRefresh string) error
}

type Service struct {
	rtRepo          RefreshTokenRepository
	encoder         AccessTokenEncoder
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

func NewService(repo RefreshTokenRepository, enc AccessTokenEncoder) AuthService {
	return &Service{
		rtRepo:          repo,
		encoder:         enc,
		accessTokenTTL:  5 * time.Minute,
		refreshTokenTTL: 30 * 24 * time.Hour,
	}
}

// -----------------------------------------------------------------------------
// GenerateForUser
// wird von SignUp / SignIn aufgerufen, um Access + Refresh Token zu erzeugen.
// -----------------------------------------------------------------------------
func (s *Service) GenerateForUser(ctx context.Context, userID uuid.UUID) (*Tokens, error) {
	// 1) AccessToken über Encoder (z.B. JWT in infra)
	accessValue, accessExp, err := s.encoder.Generate(userID, s.accessTokenTTL)
	if err != nil {
		return nil, err
	}
	access := domain.NewAccessToken(accessValue, accessExp)

	// 2) RefreshToken (raw + Hash)
	rawRefresh, hash, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := domain.NewRefreshToken(userID, hash, s.refreshTokenTTL)

	if err := s.rtRepo.Create(ctx, rt); err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  access,
		RefreshToken: rawRefresh,
	}, nil
}

// -----------------------------------------------------------------------------
// Refresh
// nimmt ein raw RefreshToken (aus Cookie), prüft es, rotiert es und gibt neue
// Tokens + userID zurück.
// -----------------------------------------------------------------------------
func (s *Service) Refresh(ctx context.Context, rawRefresh string) (*Tokens, uuid.UUID, error) {
	if rawRefresh == "" {
		return nil, uuid.Nil, ErrInvalidRefreshToken
	}

	hash := hashToken(rawRefresh)

	rt, err := s.rtRepo.FindByHash(ctx, hash)
	if err != nil || rt == nil || !rt.IsActive() {
		return nil, uuid.Nil, ErrInvalidRefreshToken
	}

	// altes Token invalidieren (Rotation)
	rt.Revoke()
	if err := s.rtRepo.Update(ctx, rt); err != nil {
		return nil, uuid.Nil, err
	}

	// neues Token-Paar ausstellen
	tokens, err := s.GenerateForUser(ctx, rt.UserID)
	if err != nil {
		return nil, uuid.Nil, err
	}

	return tokens, rt.UserID, nil
}

// -----------------------------------------------------------------------------
// Revoke
// optional: z.B. für "Logout überall".
// -----------------------------------------------------------------------------
func (s *Service) Revoke(ctx context.Context, rawRefresh string) error {
	if rawRefresh == "" {
		return ErrInvalidRefreshToken
	}

	hash := hashToken(rawRefresh)

	rt, err := s.rtRepo.FindByHash(ctx, hash)
	if err != nil || rt == nil || rt.Revoked {
		return ErrInvalidRefreshToken
	}

	rt.Revoke()
	return s.rtRepo.Update(ctx, rt)
}

// -----------------------------------------------------------------------------
// interne Helpers
// -----------------------------------------------------------------------------

func (s *Service) generateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(b)
	hash = hashToken(raw)

	return raw, hash, nil
}

func hashToken(v string) string {
	sum := sha256.Sum256([]byte(v))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
