package jwt

import (
	"errors"
	"time"

	"blogThree/internal/auth/app"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrExpiredToken            = errors.New("expired token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidClaims           = errors.New("invalid claims")
	ErrMissingSubject          = errors.New("missing subject")
	ErrInvalidSubjectFormat    = errors.New("invalid subject format")
)

type JWTEncoder struct {
	secret []byte
}

var _ app.AccessTokenEncoder = (*JWTEncoder)(nil)

func New(secret []byte) *JWTEncoder {
	return &JWTEncoder{secret: secret}
}

func (e *JWTEncoder) Generate(userID uuid.UUID, ttl time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(e.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (e *JWTEncoder) Verify(token string) (uuid.UUID, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return e.secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidClaims
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return uuid.Nil, ErrMissingSubject
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, ErrInvalidSubjectFormat
	}

	return id, nil
}
