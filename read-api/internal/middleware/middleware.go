package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

// InternalAuth rejects requests that don't carry the correct X-Internal-Secret header.
func InternalAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Internal-Secret") != secret {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// InjectUserContext reads X-User-ID and X-User-Role headers injected by the gateway.
func InjectUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if id := r.Header.Get("X-User-ID"); id != "" {
			ctx = context.WithValue(ctx, userIDKey, id)
		}
		if role := r.Header.Get("X-User-Role"); role != "" {
			ctx = context.WithValue(ctx, roleKey, role)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the authenticated user ID from context (may be empty for anonymous).
func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}

// GetRole extracts the authenticated user role from context.
func GetRole(ctx context.Context) string {
	v, _ := ctx.Value(roleKey).(string)
	return v
}
