package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestEvaluateResetTarget verifies the pure refusal rules that protect every
// database other than the configured disposable onboarding database. These
// cases are exercised WITHOUT any database connection, so no protected database
// is ever targeted.
func TestEvaluateResetTarget(t *testing.T) {
	const expected = DefaultOnboardingDatabase // db_wholesome_onboarding_test

	cases := []struct {
		name    string
		current string
		wantErr error
	}{
		{"protected frozen db refused", "db_wholesome", ErrResetTargetProtected},
		{"protected e2 test db refused", "db_wholesome_e2_test", ErrResetTargetProtected},
		{"protected e1 test db refused", "db_wholesome_e1_test", ErrResetTargetProtected},
		{"protected check is case-insensitive", "DB_WHOLESOME", ErrResetTargetProtected},
		{"unknown empty database refused", "", ErrResetTargetUnknown},
		{"whitespace-only database refused", "   ", ErrResetTargetUnknown},
		{"different database refused", "some_other_database", ErrResetTargetMismatch},
		{"configured onboarding db allowed", expected, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := evaluateResetTarget(tc.current, expected)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("evaluateResetTarget(%q) = %v, want nil", tc.current, err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("evaluateResetTarget(%q) = %v, want %v", tc.current, err, tc.wantErr)
			}
		})
	}
}

// TestDevResetDisabledByDefault ensures the feature is opt-in: with the env var
// unset the whole capability reports disabled.
func TestDevResetDisabledByDefault(t *testing.T) {
	t.Setenv("DEV_RESET_ENABLED", "")
	if DevResetEnabled() {
		t.Fatalf("DevResetEnabled() = true with DEV_RESET_ENABLED unset, want false")
	}
	t.Setenv("DEV_RESET_ENABLED", "true")
	if !DevResetEnabled() {
		t.Fatalf("DevResetEnabled() = false with DEV_RESET_ENABLED=true, want true")
	}
}

// TestDevResetPasswordDefaultAndOverride verifies the sign-in-free reset
// password is configurable and falls back to the documented default.
func TestDevResetPasswordDefaultAndOverride(t *testing.T) {
	t.Setenv("DEV_RESET_PASSWORD", "")
	if got := DevResetPassword(); got != DefaultDevResetPassword {
		t.Fatalf("DevResetPassword() = %q with no override, want %q", got, DefaultDevResetPassword)
	}
	t.Setenv("DEV_RESET_PASSWORD", "s3cret")
	if got := DevResetPassword(); got != "s3cret" {
		t.Fatalf("DevResetPassword() = %q with override, want %q", got, "s3cret")
	}
}

func newResetContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/setup/dev/reset-onboarding", strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, rec
}

// TestResetWithPasswordRejections covers the guard paths of the sign-in-free
// route. None of these reach the database: the request is refused before any
// database work begins.
func TestResetWithPasswordRejections(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("DEV_RESET_ENABLED", "true")
	t.Setenv("DEV_RESET_PASSWORD", "12345")

	cases := []struct {
		name string
		body string
		want int
	}{
		{"wrong password refused", `{"confirmation":"RESET EMS TEST SYSTEM","password":"nope"}`, http.StatusForbidden},
		{"missing password refused", `{"confirmation":"RESET EMS TEST SYSTEM"}`, http.StatusForbidden},
		{"wrong phrase refused before password", `{"confirmation":"WRONG","password":"12345"}`, http.StatusBadRequest},
		{"malformed body refused", `not json`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, rec := newResetContext(tc.body)
			ResetOnboardingDatabaseWithPassword(ctx)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tc.want, rec.Body.String())
			}
		})
	}
}

// TestResetWithPasswordDisabledByFlag ensures the route is inert unless the
// feature is explicitly enabled.
func TestResetWithPasswordDisabledByFlag(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("DEV_RESET_ENABLED", "false")

	ctx, rec := newResetContext(`{"confirmation":"RESET EMS TEST SYSTEM","password":"12345"}`)
	ResetOnboardingDatabaseWithPassword(ctx)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
