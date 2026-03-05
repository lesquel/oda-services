package jwt_test

import (
	"testing"
	"time"

	jwtutil "github.com/lesquel/oda-shared/jwt"
)

const testSecret = "test-secret-key-123"

func TestGenerate_Parse_RoundTrip(t *testing.T) {
	token, err := jwtutil.Generate("user-123", "user", testSecret, 5*time.Minute)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := jwtutil.Parse(token, testSecret)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("expected UserID=user-123, got %s", claims.UserID)
	}
	if claims.Role != "user" {
		t.Errorf("expected Role=user, got %s", claims.Role)
	}
}

func TestParse_InvalidSecret(t *testing.T) {
	token, _ := jwtutil.Generate("user-123", "user", testSecret, 5*time.Minute)
	_, err := jwtutil.Parse(token, "wrong-secret")
	if err == nil {
		t.Error("expected error when parsing with wrong secret, got nil")
	}
}

func TestParse_ExpiredToken(t *testing.T) {
	token, err := jwtutil.Generate("user-123", "user", testSecret, -1*time.Second)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	_, err = jwtutil.Parse(token, testSecret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestExtractFromHeader_Valid(t *testing.T) {
	token, _ := jwtutil.Generate("user-123", "user", testSecret, 5*time.Minute)
	extracted, err := jwtutil.ExtractFromHeader("Bearer " + token)
	if err != nil {
		t.Fatalf("ExtractFromHeader failed: %v", err)
	}
	if extracted != token {
		t.Errorf("expected extracted token to equal original token")
	}
}

func TestExtractFromHeader_Invalid(t *testing.T) {
	cases := []string{"", "token-only", "Basic abc123"}
	for _, c := range cases {
		_, err := jwtutil.ExtractFromHeader(c)
		if err == nil {
			t.Errorf("expected error for header %q, got nil", c)
		}
	}
}
