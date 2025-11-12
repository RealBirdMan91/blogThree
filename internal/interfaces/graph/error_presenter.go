// internal/interfaces/graph/error_presenter.go
package graph

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"

	apperr "blogThree/internal/errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	path := graphql.GetPath(ctx)

	var ae apperr.Error
	if errors.As(e, &ae) {
		// Nur TECHNICAL/UNKNOWN ausführlich ins Terminal loggen
		if ae.Kind() == apperr.KindTechnical || ae.Kind() == apperr.KindUnknown {
			log.Printf(
				"[TECHNICAL] code=%s kind=%s path=%v msg=%s\n%s",
				ae.Code(), ae.Kind(), path, ae.Error(), formatErrorChain(e),
			)
		}

		// Maskierter Fehler zum Client
		return &gqlerror.Error{
			Message: ae.SafeMessage(),
			Path:    path,
			Extensions: map[string]any{
				"code": string(ae.Code()),
				"kind": string(ae.Kind()),
			},
		}
	}

	// Unerwarteter Nicht-AppError: immer laut loggen inkl. Stack
	log.Printf("[UNEXPECTED] path=%v err=%v\nstack:\n%s", path, e, string(debug.Stack()))
	return &gqlerror.Error{
		Message: "internal error",
		Path:    path,
		Extensions: map[string]any{
			"code": "INTERNAL_ERROR",
			"kind": "UNKNOWN",
		},
	}
}

func RecoverFunc(ctx context.Context, rec any) error {
	// Panics immer laut loggen
	log.Printf("[PANIC] path=%v panic=%v\nstack:\n%s", graphql.GetPath(ctx), rec, string(debug.Stack()))
	// Maskiert zum Client
	return apperr.Unknown(nil)
}

func formatErrorChain(err error) string {
	var b strings.Builder
	i := 0
	for err != nil {
		if i == 0 {
			fmt.Fprintf(&b, "cause[%d]: %T: %v\n", i, err, err)
		} else {
			fmt.Fprintf(&b, " -> cause[%d]: %T: %v\n", i, err, err)
		}
		err = errors.Unwrap(err)
		i++
	}
	return b.String()
}
