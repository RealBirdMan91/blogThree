package resolvers

import (
	"blogThree/internal/interfaces/authctx"
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Auth implementiert die @auth Directive.
// Sie läuft VOR dem eigentlichen Field-Resolver.
func (r *Resolver) Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	if _, ok := authctx.UserID(ctx); !ok {
		return nil, gqlerror.Errorf("unauthorized")
	}

	return next(ctx)
}
