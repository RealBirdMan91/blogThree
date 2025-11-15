package authctx

import (
	apperr "blogThree/internal/errors"
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "auth-user-id"

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func RequireUserID(ctx context.Context) (uuid.UUID, error) {
	id, ok := UserID(ctx)
	if !ok {
		return uuid.Nil, apperr.Security("UNAUTHENTICATED", "unauthorized", nil)
	}
	return id, nil
}
