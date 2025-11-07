package security

import (
	"blogThree/internal/user/app"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{ cost int }

var _ app.PasswordHasher = (*BcryptHasher)(nil)

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *BcryptHasher) Compare(hash, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
}
