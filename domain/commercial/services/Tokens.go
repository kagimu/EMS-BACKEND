package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Token purposes. A setup token can only authorise an installation setup; a
// subscription token can only authorise a software subscription.
const (
	TokenPurposeSetup        = "setup"
	TokenPurposeSubscription = "subscription"
)

// Token statuses (explicit lifecycle).
const (
	TokenStatusIssued    = "ISSUED"
	TokenStatusValidated = "VALIDATED"
	TokenStatusRedeemed  = "REDEEMED"
	TokenStatusExpired   = "EXPIRED"
	TokenStatusRevoked   = "REVOKED"
)

// GenerateRawToken returns a cryptographically secure, URL-safe 256-bit token.
// The raw value is shown once to the client and is never persisted.
func GenerateRawToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashToken returns the SHA-256 hex digest used for at-rest storage and lookup.
// Tokens are looked up by hash, so a token can never be matched to another
// order.
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
