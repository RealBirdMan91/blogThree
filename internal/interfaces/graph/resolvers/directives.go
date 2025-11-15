package resolvers

import (
	apperr "blogThree/internal/errors"
	"blogThree/internal/interfaces/authctx"
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// Auth implementiert die @auth Directive.
// Sie läuft VOR dem eigentlichen Field-Resolver.
func (r *Resolver) Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	if _, ok := authctx.UserID(ctx); !ok {
		return nil, apperr.Security("UNAUTHENTICATED", "unauthorized", nil)
	}

	return next(ctx)
}
