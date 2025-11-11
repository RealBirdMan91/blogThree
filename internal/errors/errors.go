package errors

import "errors"

type Kind string

const (
	KindDomain     Kind = "DOMAIN"
	KindValidation Kind = "VALIDATION"
	KindTechnical  Kind = "TECHNICAL"
	KindSecurity   Kind = "SECURITY"
	KindUnknown    Kind = "UNKNOWN"
)

type Code string

type Error interface {
	error
	Code() Code
	Kind() Kind
	SafeMessage() string
	Extensions() map[string]any
	Unwrap() error
}

type appError struct {
	kind    Kind
	code    Code
	msg     string
	safe    string
	ext     map[string]any
	wrapped error
}

func (e *appError) Error() string       { return e.msg }
func (e *appError) Code() Code          { return e.code }
func (e *appError) Kind() Kind          { return e.kind }
func (e *appError) SafeMessage() string { return e.safe }
func (e *appError) Extensions() map[string]any {
	if e.ext == nil {
		return map[string]any{}
	}
	return e.ext
}
func (e *appError) Unwrap() error { return e.wrapped }

func New(kind Kind, code Code, safe string, ext map[string]any) Error {
	return &appError{
		kind: kind,
		code: code,
		safe: safe,
		msg:  string(code) + ": " + safe,
		ext:  ext,
	}
}

func Wrap(kind Kind, code Code, safe string, ext map[string]any, cause error) Error {
	if cause == nil {
		return New(kind, code, safe, ext)
	}
	return &appError{
		kind:    kind,
		code:    code,
		safe:    safe,
		msg:     string(code) + ": " + cause.Error(),
		ext:     ext,
		wrapped: cause,
	}
}

func Domain(code Code, safe string, ext map[string]any) Error {
	return New(KindDomain, code, safe, ext)
}

func Validation(code Code, safe string, ext map[string]any) Error {
	return New(KindValidation, code, safe, ext)
}

func Security(code Code, safe string, ext map[string]any) Error {
	return New(KindSecurity, code, safe, ext)
}

func Technical(code Code, safe string, ext map[string]any, cause error) Error {
	return Wrap(KindTechnical, code, safe, ext, cause)
}

func Unknown(cause error) Error {
	if cause == nil {
		cause = errors.New("unknown error")
	}
	return Wrap(KindUnknown, "INTERNAL_ERROR", "internal error", nil, cause)
}
