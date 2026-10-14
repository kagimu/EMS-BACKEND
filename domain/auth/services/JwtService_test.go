package services

import "testing"

// TestGenerateAndValidateTokenRoundTrip freezes the JWT contract that the
// existing Education authentication flow depends on: a token generated for a
// user number must validate and carry the same user_number claim.
func TestGenerateAndValidateTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-secret")

	svc := JWTAuthService()

	token, err := svc.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned an empty token")
	}

	parsed, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected the generated token to be valid")
	}

	claims, ok := parsed.Claims.(*AuthCustomClaims)
	if !ok {
		t.Fatalf("unexpected claims type %T", parsed.Claims)
	}
	if claims.UserNumber != "user-123" {
		t.Fatalf("UserNumber = %q, want %q", claims.UserNumber, "user-123")
	}
}

// TestValidateTokenRejectsDifferentSecret proves that validation actually
// depends on the configured secret and is not silently accepting anything.
func TestValidateTokenRejectsDifferentSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret-a")
	token, err := JWTAuthService().GenerateToken("user-1")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	t.Setenv("JWT_SECRET", "secret-b")
	if _, err := JWTAuthService().ValidateToken(token); err == nil {
		t.Fatal("expected validation to fail when the signing secret differs")
	}
}

// TestGetSecretKeyHasNoFallback is the Release 1 security regression guard for
// fix B1: a missing JWT_SECRET must never produce a predictable signing key.
func TestGetSecretKeyHasNoFallback(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	if got := getSecretKey(); got != "" {
		t.Fatalf("getSecretKey() = %q, want empty string when JWT_SECRET is unset", got)
	}
}
