package middlewares

import (
	"context"
	"net/http"

	"github.com/govnocods/RedChat/internal/auth"
)

func (m *Middlewares) APIAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}

		tokenStr := cookie.Value

		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized: invalid token"}`))
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserId)
		ctx = context.WithValue(ctx, usernameKey, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
