package config

import (
	"log"
	"os"
)

const defaultSecret = "change-internal-secret-in-production"

type Config struct {
	DatabaseURL    string
	Port           string
	InternalSecret string
}

func Load() *Config {
	secret := getEnv("INTERNAL_SECRET", defaultSecret)
	if secret == defaultSecret {
		log.Println("WARNING: using default INTERNAL_SECRET – change before deploying")
	}

	// DATABASE_URL_READ allows pointing to the read replica in dev.
	// Falls back to DATABASE_URL (which defaults to primary on :5432).
	dbURL := getEnv("DATABASE_URL_READ", "")
	if dbURL == "" {
		dbURL = getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/oda?sslmode=disable")
	}

	return &Config{
		DatabaseURL:    dbURL,
		Port:           getEnv("PORT", "8083"),
		InternalSecret: secret,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
