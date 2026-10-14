// Package security contains source-level regression tests for the EMS Release 1
// security fixes. These tests do not require a database or a running server;
// they assert that the fixed defects cannot be reintroduced unnoticed.
package security

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// backendFile resolves a path relative to the backend module root and returns
// both the resolved path and its contents.
func backendFile(t *testing.T, rel string) (string, string) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve the test file path")
	}

	backendRoot := filepath.Dir(filepath.Dir(thisFile))
	path := filepath.Join(backendRoot, filepath.FromSlash(rel))

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	return path, string(content)
}

// TestJwtServiceHasNoFallbackSecret guards Release 1 fix B1.
func TestJwtServiceHasNoFallbackSecret(t *testing.T) {
	path, src := backendFile(t, "domain/auth/services/JwtService.go")

	if strings.Contains(src, `secret = "secret"`) {
		t.Errorf("JWT fallback secret is present again in %s", path)
	}
	if !strings.Contains(src, `os.Getenv("JWT_SECRET")`) {
		t.Errorf("%s must read JWT_SECRET from the environment", path)
	}
}

// TestStartupRequiresJwtSecret guards the startup half of Release 1 fix B1.
func TestStartupRequiresJwtSecret(t *testing.T) {
	path, src := backendFile(t, "main.go")

	if !strings.Contains(src, `os.Getenv("JWT_SECRET") == ""`) {
		t.Errorf("%s must refuse to start when JWT_SECRET is empty", path)
	}
}

// TestEmailSenderReadsCredentialsFromEnv guards Release 1 fix B2.
func TestEmailSenderReadsCredentialsFromEnv(t *testing.T) {
	path, src := backendFile(t, "domain/commons/emailSender/EmailSender.go")

	if strings.Contains(src, "andysteve2222") {
		t.Errorf("hardcoded sender address is present again in %s", path)
	}
	for _, required := range []string{
		`os.Getenv("MAIL_USERNAME")`,
		`os.Getenv("MAIL_PASSWORD")`,
		`os.Getenv("MAIL_HOST")`,
		`os.Getenv("MAIL_PORT")`,
	} {
		if !strings.Contains(src, required) {
			t.Errorf("%s must read %s from the environment", path, required)
		}
	}
}

// TestCorsDoesNotAcceptEveryOrigin guards Release 1 fix B3.
func TestCorsDoesNotAcceptEveryOrigin(t *testing.T) {
	path, src := backendFile(t, "api/routes/Routes.go")

	if strings.Contains(src, "AllowOriginFunc") {
		t.Errorf("permissive origin callback is present again in %s", path)
	}
	if !strings.Contains(src, `os.Getenv("APP_CORS_ORIGINS")`) {
		t.Errorf("%s must allow CORS origins to be configured via APP_CORS_ORIGINS", path)
	}
}
