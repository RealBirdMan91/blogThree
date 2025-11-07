package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	email     Email
	password  PasswordHash
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email Email, passwordHash PasswordHash) (*User, error) {
	now := time.Now().UTC()

	return &User{
		id:        uuid.New(),
		email:     email,
		password:  passwordHash,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RehydrateUser(
	id uuid.UUID,
	email Email,
	password PasswordHash,
	createdAt, updatedAt time.Time,
) (*User, error) {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (u *User) ID() uuid.UUID              { return u.id }
func (u *User) Email() Email               { return u.email }
func (u *User) PasswordHash() PasswordHash { return u.password }
func (u *User) CreatedAt() time.Time       { return u.createdAt }
func (u *User) UpdatedAt() time.Time       { return u.updatedAt }
