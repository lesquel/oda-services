package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtutil "github.com/lesquel/oda-shared/jwt"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

// OptionalJWTAuth validates a JWT if present, injecting user context. Never rejects requests.
func OptionalJWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenStr, err := jwtutil.ExtractFromHeader(authHeader)
				if err == nil {
					if claims, err := jwtutil.Parse(tokenStr, jwtSecret); err == nil {
						ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
						ctx = context.WithValue(ctx, roleKey, claims.Role)
						r = r.WithContext(ctx)
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireJWTAuth validates a JWT, returning 401 if missing or invalid.
func RequireJWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			tokenStr, err := jwtutil.ExtractFromHeader(authHeader)
			if err != nil {
				http.Error(w, `{"error":"invalid authorization header"}`, http.StatusUnauthorized)
				return
			}
			claims, err := jwtutil.Parse(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdminAuth requires a valid JWT with role=admin.
func RequireAdminAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequireJWTAuth(jwtSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(roleKey).(string)
			if !strings.EqualFold(role, "admin") {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

// GetUserID extracts user ID from context.
func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}

// GetRole extracts user role from context.
func GetRole(ctx context.Context) string {
	v, _ := ctx.Value(roleKey).(string)
	return v
}
