package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims contains the JWT payload.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Generate creates a signed JWT for the given userID and role.
// accessTTL controls the expiry; pass 0 to use the default of 15 minutes.
func Generate(userID, role, secret string, accessTTL time.Duration) (string, error) {
	if accessTTL == 0 {
		accessTTL = 15 * time.Minute
	}
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefresh creates a long-lived refresh token (30 days).
func GenerateRefresh(userID, role, secret string) (string, error) {
	return Generate(userID, role, secret, 30*24*time.Hour)
}

// Parse validates and parses a JWT string, returning its Claims.
func Parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ExtractFromHeader strips the "Bearer " prefix from an Authorization header value.
func ExtractFromHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}
