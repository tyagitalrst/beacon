package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	Env         string // "dev" | "prod"
}

func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://beacon:beacon@localhost:5435/beacon?sslmode=disable"),
		Env:         getenv("APP_ENV", "dev"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
