package app

import (
	apperr "blogThree/internal/errors"
)

const (
	CodeInvalidTitle apperr.Code = "POST_INVALID_TITLE"
	CodeInvalidBody  apperr.Code = "POST_INVALID_BODY"

	CodePostNotFound      apperr.Code = "POST_NOT_FOUND"
	CodeAuthorNotFound    apperr.Code = "POST_AUTHOR_NOT_FOUND"
	CodePostPersistFailed apperr.Code = "POST_PERSIST_FAILED"
	CodePostSelectFailed  apperr.Code = "POST_SELECT_FAILED"
	CodePostsListFailed   apperr.Code = "POSTS_LIST_FAILED"
	CodeSlugAlreadyExists apperr.Code = "POST_SLUG_ALREADY_EXISTS"
	CodePostUnknownError  apperr.Code = "POST_UNKNOWN_ERROR"
	CodeAuthorCheckFailed apperr.Code = "AUTHOR_CHECK_FAILED"
)

var (
	ErrPostNotFound   = apperr.Domain(CodePostNotFound, "post not found", nil)
	ErrAuthorNotFound = apperr.Domain(CodeAuthorNotFound, "author not found", nil)
)

func NewInvalidTitleError(reason string) apperr.Error {
	ext := map[string]any{"field": "title"}
	if reason != "" {
		ext["reason"] = reason
	}
	return apperr.Validation(CodeInvalidTitle, "invalid title", ext)
}

func NewInvalidBodyError(reason string) apperr.Error {
	return apperr.Validation(CodeInvalidBody, "invalid body", map[string]any{"field": "body"})
}

func NewPostPersistFailed(cause error) apperr.Error {
	return apperr.Technical(CodePostPersistFailed, "could not create post", nil, cause)
}
func NewPostSelectFailed(cause error) apperr.Error {
	return apperr.Technical(CodePostSelectFailed, "could not load post", nil, cause)
}
func NewPostsListFailed(cause error) apperr.Error {
	return apperr.Technical(CodePostsListFailed, "could not list posts", nil, cause)
}

func NewUnknownPostError(cause error) apperr.Error {
	return apperr.Unknown(cause)
}
func NewAuthorCheckFailed(cause error) apperr.Error {
	return apperr.Technical(CodeAuthorCheckFailed, "could not verify author", nil, cause)
}
