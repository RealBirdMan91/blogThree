package jwt

import (
	"time"

	"blogThree/internal/auth/app"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
