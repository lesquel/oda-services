package hasher_test

import (
	"testing"

	"github.com/lesquel/oda-shared/hasher"
)

func TestHashPassword_CheckPassword(t *testing.T) {
	password := "s3cr3tP@ssword"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == password {
		t.Error("hash should not equal plain password")
	}
	if !hasher.CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password")
	}
	if hasher.CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "same-password"
	hash1, _ := hasher.HashPassword(password)
	hash2, _ := hasher.HashPassword(password)
	if hash1 == hash2 {
		t.Error("bcrypt should produce different hashes each time (different salts)")
	}
}
