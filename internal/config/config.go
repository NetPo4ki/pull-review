package config

import "os"

type Config struct {
	DBDSN    string
	HTTPPort string
	AppEnv   string
	LogLevel string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		DBDSN:    getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		HTTPPort: getenv("HTTP_PORT", "8080"),
		AppEnv:   getenv("APP_ENV", "local"),
		LogLevel: getenv("LOG_LEVEL", "info"),
	}
}
