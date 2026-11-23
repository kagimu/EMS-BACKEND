package providers

import (
	"context"
	"errors"
	"testing"
)

func TestAirtelRemainsDisabledWithoutVerifiedContract(t *testing.T) {
	provider := &Airtel{cfg: AirtelConfig{
		Enabled: true, ContractVerified: false,
		BaseURL: "https://invalid.example", ClientID: "id", ClientSecret: "secret",
	}}
	if provider.Enabled() {
		t.Fatal("Airtel must not be enabled without a verified contract flag")
	}
	_, err := provider.Initiate(context.Background(), Request{})
	if err == nil {
		t.Fatal("Airtel initiation must fail safely when its contract is unverified")
	}
	if !errors.Is(err, ErrContract) && !errors.Is(err, ErrDisabled) {
		t.Fatalf("unexpected Airtel disabled error: %v", err)
	}
}
