package app

import (
	apperr "blogThree/internal/errors"
)

// -------------------- ERROR CODES --------------------

const (
	CodeNoRequestInContext apperr.Code = "AUTH_NO_REQUEST_IN_CONTEXT"
	CodeNoRefreshToken     apperr.Code = "AUTH_NO_REFRESH_TOKEN"
)

// -------------------- CONSTRUCTORS --------------------

// Technischer Fehler: der HTTP-Request hängt nicht im Context.
// Das ist eher ein Infrastrukturproblem (Middleware vergessen, etc.).
func NewNoRequestInContextError() apperr.Error {
	return apperr.Technical(
		CodeNoRequestInContext,
		"internal error",
		nil,
		nil,
	)
}

// Security-Fehler: es gibt keinen (oder einen leeren) Refresh-Token.
// Das ist ein Auth-Problem -> KindSecurity.
func NewNoRefreshTokenError() apperr.Error {
	return apperr.Security(
		CodeNoRefreshToken,
		"not authenticated",
		nil,
	)
}
