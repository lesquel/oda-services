package hasher

import "golang.org/x/crypto/bcrypt"

const cost = 12

// HashPassword hashes a plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}

// CheckPassword returns true when the plain-text password matches the hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
