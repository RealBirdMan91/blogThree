package app

import (
	"errors"

	apperr "blogThree/internal/errors"
)

const (
	CodeInvalidEmail       apperr.Code = "INVALID_EMAIL"
	CodeWeakPassword       apperr.Code = "WEAK_PASSWORD"
	CodeEmailAlreadyExists apperr.Code = "USER_EMAIL_ALREADY_EXISTS"
	CodeInvalidCredentials apperr.Code = "AUTH_INVALID_CREDENTIALS"
	CodeUserNotFound       apperr.Code = "USER_NOT_FOUND"

	CodeUserExistsCheckFailed apperr.Code = "USER_EXISTS_CHECK_FAILED"
	CodeUserInsertFailed      apperr.Code = "USER_INSERT_FAILED"
	CodeUserSelectFailed      apperr.Code = "USER_SELECT_FAILED"
	CodeUserListFailed        apperr.Code = "USER_LIST_FAILED"
	CodePasswordHashFailed    apperr.Code = "PASSWORD_HASH_FAILED"
)

var (
	ErrEmailAlreadyExists = apperr.Domain(
		CodeEmailAlreadyExists,
		"email already in use",
		nil,
	)

	ErrInvalidCredentials = apperr.Security(
		CodeInvalidCredentials,
		"invalid email or password",
		nil,
	)

	ErrUserNotFound = apperr.Domain(
		CodeUserNotFound,
		"user not found",
		nil,
	)
)

func NewInvalidEmailError() apperr.Error {
	return apperr.Validation(
		CodeInvalidEmail,
		"invalid email",
		map[string]any{"field": "email"},
	)
}

func NewWeakPasswordError(reason string) apperr.Error {
	ext := map[string]any{"field": "password"}
	if reason != "" {
		ext["reason"] = reason
	}
	return apperr.Validation(
		CodeWeakPassword,
		"password does not meet complexity requirements",
		ext,
	)
}

// Technische Wrapper:

func NewUserExistsCheckFailed(cause error) apperr.Error {
	return apperr.Technical(
		CodeUserExistsCheckFailed,
		"could not check if user exists",
		nil,
		cause,
	)
}

func NewUserInsertFailed(cause error) apperr.Error {
	return apperr.Technical(
		CodeUserInsertFailed,
		"could not create user",
		nil,
		cause,
	)
}

func NewUserSelectFailed(cause error) apperr.Error {
	return apperr.Technical(
		CodeUserSelectFailed,
		"could not load user",
		nil,
		cause,
	)
}

func NewUserListFailed(cause error) apperr.Error {
	return apperr.Technical(
		CodeUserListFailed,
		"could not list users",
		nil,
		cause,
	)
}

func NewPasswordHashFailed(cause error) apperr.Error {
	return apperr.Technical(
		CodePasswordHashFailed,
		"could not hash password",
		nil,
		cause,
	)
}

// Helper um zu prüfen, ob ein Error unser Typ ist
func IsAppError(err error) bool {
	var ae apperr.Error
	return errors.As(err, &ae)
}
