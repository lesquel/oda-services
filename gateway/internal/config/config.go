package config

import (
	"log"
	"os"
)

const defaultSecret = "change-internal-secret-in-production"

type Config struct {
	Port           string
	JWTSecret      string
	InternalSecret string
	WriteAPIURL    string
	ReadAPIURL     string
}

func Load() *Config {
	secret := getEnv("INTERNAL_SECRET", defaultSecret)
	if secret == defaultSecret {
		log.Println("WARNING: using default INTERNAL_SECRET – change before deploying")
	}
	return &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "change-jwt-secret-in-production"),
		InternalSecret: secret,
		WriteAPIURL:    getEnv("WRITE_API_URL", "http://localhost:8082"),
		ReadAPIURL:     getEnv("READ_API_URL", "http://localhost:8083"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
