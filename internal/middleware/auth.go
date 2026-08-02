package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	pkgauth "blog-api/pkg/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func respondJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondJSON(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondJSON(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				respondJSON(w, http.StatusUnauthorized, "empty token")
				return
			}

			claims, err := pkgauth.ParseToken(tokenString, secret)
			if err != nil {
				respondJSON(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
