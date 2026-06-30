package config

import (
	"errors"
	"os"
	"strings"
)

const defaultPort = "8080"

type Config struct {
	AppEnv      string
	Port        string
	OpenAPIPath string
	DatabaseURL string
	UploadDir   string
	CORSOrigins []string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      envOrDefault("APP_ENV", "local"),
		Port:        envOrDefault("PORT", defaultPort),
		OpenAPIPath: resolveOpenAPIPath(os.Getenv("OPENAPI_PATH")),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://aromatype:aromatype@localhost:5432/aromatype?sslmode=disable"),
		UploadDir:   envOrDefault("UPLOAD_DIR", "uploads"),
		CORSOrigins: splitCSV(envOrDefault("CORS_ALLOWED_ORIGINS", "*")),
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
	if cfg.UploadDir == "" {
		return Config{}, errors.New("UPLOAD_DIR must not be empty")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
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
