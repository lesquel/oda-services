package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lesquel/oda-read-api/internal/middleware"
)

// ── Tests for InternalAuth ────────────────────────────────────────────────────────

const testInternalSecret = "test-secret-key"

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
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
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

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
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

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

// ── Tests for InjectUserContext ─────────────────────────────────────────────────

func TestInjectUserContext_WithUserID(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		if userID != "user-123" {
			t.Errorf("expected userID 'user-123', got '%s'", userID)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.InjectUserContext(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestInjectUserContext_WithRole(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := middleware.GetRole(r.Context())
		if role != "admin" {
			t.Errorf("expected role 'admin', got '%s'", role)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.InjectUserContext(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-Role", "admin")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestInjectUserContext_WithBothHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		role := middleware.GetRole(r.Context())
		if userID != "user-456" || role != "moderator" {
			t.Errorf("expected userID 'user-456' and role 'moderator', got '%s' and '%s'", userID, role)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.InjectUserContext(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-456")
	req.Header.Set("X-User-Role", "moderator")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestInjectUserContext_MissingHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		role := middleware.GetRole(r.Context())
		if userID != "" {
			t.Errorf("expected empty userID, got '%s'", userID)
		}
		if role != "" {
			t.Errorf("expected empty role, got '%s'", role)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.InjectUserContext(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

// TestGetUserID_FromContext - This test verifies that InjectUserContext correctly
// sets the user ID in context. The actual key is private, so we test through
// the middleware behavior in TestInjectUserContext_WithUserID instead.
func TestGetUserID_FromContext(t *testing.T) {
	// Skip this test - the context key is private and can't be tested directly.
	// The InjectUserContext tests above verify the behavior.
	t.Skip("context key is private - behavior tested via middleware tests")
}

func TestGetUserID_EmptyContext(t *testing.T) {
	ctx := context.Background()
	userID := middleware.GetUserID(ctx)

	if userID != "" {
		t.Errorf("expected empty string, got '%s'", userID)
	}
}

// TestGetRole_FromContext - Same reason as above.
func TestGetRole_FromContext(t *testing.T) {
	t.Skip("context key is private - behavior tested via middleware tests")
}

// ── Table-Driven Tests ────────────────────────────────────────────────────────────

type authTestCase struct {
	name           string
	secret         string
	headerSecret   string
	expectedStatus int
}

func TestInternalAuth_TableDriven(t *testing.T) {
	tests := []authTestCase{
		{
			name:           "valid secret",
			secret:         "my-secret",
			headerSecret:   "my-secret",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid secret",
			secret:         "my-secret",
			headerSecret:   "wrong-secret",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "empty secret",
			secret:         "my-secret",
			headerSecret:   "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.InternalAuth(tt.secret)(next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.headerSecret != "" {
				req.Header.Set("X-Internal-Secret", tt.headerSecret)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
