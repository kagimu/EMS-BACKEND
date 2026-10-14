package controllers

import (
	"errors"
	"testing"
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
