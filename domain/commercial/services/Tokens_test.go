package services

import "testing"

func TestGenerateRawTokenIsUniqueAndOpaque(t *testing.T) {
	first, err := GenerateRawToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	second, err := GenerateRawToken()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if first == "" || second == "" {
		t.Fatal("empty token")
	}
	if first == second {
		t.Fatal("tokens are not unique")
	}
	if len(first) < 32 {
		t.Fatalf("token too short: %q", first)
	}
}

func TestHashTokenIsDeterministicAndOpaque(t *testing.T) {
	raw := "test-token-value"
	h1 := HashToken(raw)
	h2 := HashToken(raw)
	if h1 != h2 {
		t.Fatal("hash is not deterministic")
	}
	if len(h1) != 64 {
		t.Fatalf("hash length = %d, want 64 (sha256 hex)", len(h1))
	}
	if h1 == raw {
		t.Fatal("hash must not equal the raw token")
	}
	if HashToken("other") == h1 {
		t.Fatal("different tokens must hash differently")
	}
}
