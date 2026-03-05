package config

import (
	"log"
	"os"
)

const defaultSecret = "change-internal-secret-in-production"

// Config holds write-api configuration.
type Config struct {
	DatabaseURL    string
	Port           string
	JWTSecret      string
	InternalSecret string // shared with gateway – rejects unauthenticated internal calls
	AdminEmail     string
	AdminPassword  string
	NATSURL        string
}

// Load reads configuration from environment variables.
func Load() *Config {
	secret := getEnv("INTERNAL_SECRET", defaultSecret)
	if secret == defaultSecret {
		log.Println("WARNING: using default INTERNAL_SECRET – change before deploying")
	}
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/oda?sslmode=disable"),
		Port:           getEnv("PORT", "8082"),
		JWTSecret:      getEnv("JWT_SECRET", "change-jwt-secret-in-production"),
		InternalSecret: secret,
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@oda.com"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", ""),
		NATSURL:        getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
