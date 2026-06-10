package config

import (
	"errors"
	"os"
)

const defaultPort = "8080"

type Config struct {
	AppEnv      string
	Port        string
	OpenAPIPath string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      envOrDefault("APP_ENV", "local"),
		Port:        envOrDefault("PORT", defaultPort),
		OpenAPIPath: resolveOpenAPIPath(os.Getenv("OPENAPI_PATH")),
	}

	if cfg.Port == "" {
		return Config{}, errors.New("PORT must not be empty")
	}
	if cfg.OpenAPIPath == "" {
		return Config{}, errors.New("OpenAPI spec file was not found")
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
