package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSplitCSVTrimsAndSkipsEmptyValues(t *testing.T) {
	got := splitCSV(" *, http://localhost:5173, ,https://example.com ")
	want := []string{"*", "http://localhost:5173", "https://example.com"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected csv split: got %#v want %#v", got, want)
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("AROMATYPE_TEST_ENV", "configured")

	if got := envOrDefault("AROMATYPE_TEST_ENV", "fallback"); got != "configured" {
		t.Fatalf("expected configured value, got %q", got)
	}
	if got := envOrDefault("AROMATYPE_MISSING_ENV", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback value, got %q", got)
	}
}

func TestResolveOpenAPIPathUsesConfiguredExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "openapi.yaml")
	if err := os.WriteFile(path, []byte("openapi: 3.0.3\n"), 0o644); err != nil {
		t.Fatalf("write openapi fixture: %v", err)
	}

	if got := resolveOpenAPIPath(path); got != path {
		t.Fatalf("expected configured path, got %q", got)
	}
}

func TestLoadUsesEnvironment(t *testing.T) {
	dir := t.TempDir()
	openAPIPath := filepath.Join(dir, "openapi.yaml")
	if err := os.WriteFile(openAPIPath, []byte("openapi: 3.0.3\n"), 0o644); err != nil {
		t.Fatalf("write openapi fixture: %v", err)
	}

	t.Setenv("APP_ENV", "test")
	t.Setenv("PORT", "9090")
	t.Setenv("OPENAPI_PATH", openAPIPath)
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("UPLOAD_DIR", "tmp-uploads")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://one.test, https://two.test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.AppEnv != "test" || cfg.Port != "9090" || cfg.OpenAPIPath != openAPIPath {
		t.Fatalf("unexpected config: %#v", cfg)
	}
	wantOrigins := []string{"https://one.test", "https://two.test"}
	if !reflect.DeepEqual(cfg.CORSOrigins, wantOrigins) {
		t.Fatalf("unexpected origins: got %#v want %#v", cfg.CORSOrigins, wantOrigins)
	}
}
