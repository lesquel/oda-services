// Package middleware provides internal authentication and context helpers for write-api.
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "userRole"
)

// InternalAuth rejects requests that don't carry the correct X-Internal-Secret header.
// This ensures write-api is only reachable through the gateway.
func InternalAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Internal-Secret") != secret {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// InjectUserContext reads X-User-ID and X-User-Role forwarded by the gateway
// and stores them in the request context.
func InjectUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			ctx = context.WithValue(ctx, UserIDKey, uid)
		}
		if role := r.Header.Get("X-User-Role"); role != "" {
			ctx = context.WithValue(ctx, RoleKey, role)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireUser returns 401 when X-User-ID is absent from context.
func RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(UserIDKey).(string); !ok {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin returns 403 when the user role is not "admin".
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleKey).(string)
		if role != "admin" {
			http.Error(w, `{"error":"forbidden: admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the user ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok && id != ""
}

// GetRole extracts the user role from context.
func GetRole(ctx context.Context) string {
	role, _ := ctx.Value(RoleKey).(string)
	if role == "" {
		return "user"
	}
	return role
}
