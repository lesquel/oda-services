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
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/oda?sslmode=disable"),
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
