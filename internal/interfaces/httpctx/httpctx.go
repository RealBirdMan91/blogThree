package httpctx

import (
	"context"
	"net/http"
)

type ctxKey string

const (
	responseWriterKey ctxKey = "resp-writer"
	requestKey        ctxKey = "req"
)

func Inject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		ctx = context.WithValue(ctx, requestKey, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	return w, ok
}

func Request(ctx context.Context) (*http.Request, bool) {
	r, ok := ctx.Value(requestKey).(*http.Request)
	return r, ok
}
