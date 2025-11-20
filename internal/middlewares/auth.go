package middlewares

import (
	"context"
	"net/http"

	"github.com/govnocods/RedChat/internal/auth"
)

type contextKey string

const (
	userIDKey   contextKey = "userID"
	usernameKey contextKey = "username"
)

func (m *Middlewares) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenStr := cookie.Value

		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserId)
		ctx = context.WithValue(ctx, usernameKey, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
