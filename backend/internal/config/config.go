package config

import (
	"errors"
	"os"
	"time"
)

const defaultPort = "8080"

type Config struct {
	AppEnv      string
	Port        string
	OpenAPIPath string
	DatabaseURL string
	JWTSecret   string
	JWTTTL      time.Duration
}

func Load() (Config, error) {
	jwtTTL, err := time.ParseDuration(envOrDefault("JWT_TTL", "24h"))
	if err != nil {
		return Config{}, errors.New("JWT_TTL must be a valid duration")
	}

	cfg := Config{
		AppEnv:      envOrDefault("APP_ENV", "local"),
		Port:        envOrDefault("PORT", defaultPort),
		OpenAPIPath: resolveOpenAPIPath(os.Getenv("OPENAPI_PATH")),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://aromatype:aromatype@localhost:5432/aromatype?sslmode=disable"),
		JWTSecret:   envOrDefault("JWT_SECRET", "local-dev-secret-change-me"),
		JWTTTL:      jwtTTL,
	}

	if cfg.Port == "" {
		return Config{}, errors.New("PORT must not be empty")
	}
	if cfg.OpenAPIPath == "" {
		return Config{}, errors.New("OpenAPI spec file was not found")
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL must not be empty")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET must not be empty")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func resolveOpenAPIPath(configured string) string {
	candidates := []string{
		configured,
		"../docs/api/openapi.yaml",
		"docs/api/openapi.yaml",
		"/app/docs/api/openapi.yaml",
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}
