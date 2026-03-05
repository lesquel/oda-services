package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lesquel/oda-write-api/internal/middleware"
)

const testInternalSecret = "super-internal-secret"

func TestInternalAuth_ValidSecret(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.InternalAuth(testInternalSecret)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Secret", testInternalSecret)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestInternalAuth_InvalidSecret(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.InternalAuth(testInternalSecret)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Internal-Secret", "wrong-secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInternalAuth_MissingSecret(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.InternalAuth(testInternalSecret)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestInjectUserContext_SetsValues(t *testing.T) {
	var capturedID, capturedRole string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID, _ = middleware.GetUserID(r.Context())
		capturedRole = middleware.GetRole(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.InjectUserContext(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-123")
	req.Header.Set("X-User-Role", "user")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if capturedID != "user-123" {
		t.Errorf("expected userID user-123, got %s", capturedID)
	}
	if capturedRole != "user" {
		t.Errorf("expected role user, got %s", capturedRole)
	}
}

func TestRequireUser_WithUser(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.RequireUser(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = reqWithUserID(req, "user-123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireUser_Without(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.RequireUser(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAdmin_WithAdmin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.RequireAdmin(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = reqWithRole(req, "admin")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestRequireAdmin_WithUser(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware.RequireAdmin(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = reqWithRole(req, "user")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}

// helpers: run InjectUserContext to inject context values
func reqWithUserID(r *http.Request, id string) *http.Request {
	r.Header.Set("X-User-ID", id)
	var enriched *http.Request
	next := http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
		enriched = r2
	})
	middleware.InjectUserContext(next).ServeHTTP(httptest.NewRecorder(), r)
	return enriched
}

func reqWithRole(r *http.Request, role string) *http.Request {
	r.Header.Set("X-User-ID", "user-123")
	r.Header.Set("X-User-Role", role)
	var enriched *http.Request
	next := http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
		enriched = r2
	})
	middleware.InjectUserContext(next).ServeHTTP(httptest.NewRecorder(), r)
	return enriched
}
