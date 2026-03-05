package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lesquel/oda-gateway/internal/middleware"
)

const gatewaySecret = "gateway-test-secret"

func makeToken(userID, role, secret string, valid bool) string {
	expiry := time.Now().Add(time.Hour)
	if !valid {
		expiry = time.Now().Add(-time.Hour)
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestOptionalJWTAuth_WithValidToken(t *testing.T) {
	var gotID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = middleware.GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.OptionalJWTAuth(gatewaySecret)(next)
	tokenStr := makeToken("u1", "user", gatewaySecret, true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	h.ServeHTTP(httptest.NewRecorder(), req)
	if gotID != "u1" {
		t.Errorf("expected userID u1, got %q", gotID)
	}
}

func TestOptionalJWTAuth_WithoutToken(t *testing.T) {
	var gotID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = middleware.GetUserID(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.OptionalJWTAuth(gatewaySecret)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if gotID != "" {
		t.Errorf("expected empty userID, got %q", gotID)
	}
}

func TestRequireJWTAuth_ValidToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.RequireJWTAuth(gatewaySecret)(next)
	tokenStr := makeToken("u2", "user", gatewaySecret, true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireJWTAuth_MissingToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.RequireJWTAuth(gatewaySecret)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestRequireJWTAuth_ExpiredToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.RequireJWTAuth(gatewaySecret)(next)
	tokenStr := makeToken("u3", "user", gatewaySecret, false)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAdminAuth_WithAdminRole(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.RequireAdminAuth(gatewaySecret)(next)
	tokenStr := makeToken("admin1", "admin", gatewaySecret, true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireAdminAuth_WithUserRole(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := middleware.RequireAdminAuth(gatewaySecret)(next)
	tokenStr := makeToken("normaluser", "user", gatewaySecret, true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}
