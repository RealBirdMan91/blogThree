package authctx

import (
	"net/http"
	"strings"

	authapp "blogThree/internal/auth/app"
)

func Middleware(encoder authapp.AccessTokenEncoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// anonym weiter
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				// ungültiges Format -> hier kannst du auch 401 senden
				next.ServeHTTP(w, r)
				return
			}

			token := parts[1]
			userID, err := encoder.Verify(token)
			if err != nil {
				// invalides Token -> optional 401
				// http.Error(w, "invalid token", http.StatusUnauthorized); return
				next.ServeHTTP(w, r)
				return
			}

			ctx := WithUserID(r.Context(), userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
