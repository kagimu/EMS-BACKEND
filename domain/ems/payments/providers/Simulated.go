package providers

import (
	"context"
	"os"
	"strings"
)

const Simulated = "SIMULATED"

// Simulated is a development-only provider. It never contacts a telecom
// network. It exists so the full order -> payment record -> owner-console
// verification path can be exercised without merchant credentials.
type SimulatedProvider struct{}

func NewSimulated() *SimulatedProvider { return &SimulatedProvider{} }

func (s *SimulatedProvider) Name() string { return Simulated }

func (s *SimulatedProvider) DisplayName() string { return "Simulated (test)" }

func (s *SimulatedProvider) Enabled() bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("ENVIRONMENT")), "production") {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEV_SIMULATED_PAYMENTS_ENABLED")), "true")
}

func (s *SimulatedProvider) Environment() string { return "test" }

func (s *SimulatedProvider) Initiate(_ context.Context, request Request) (Result, error) {
	return Result{
		Status:            "SUCCESSFUL",
		ProviderStatus:    "SIMULATED_SUCCESS",
		ProviderReference: request.PaymentNumber,
		ProviderAmount:    request.Amount,
		ProviderCurrency:  request.Currency,
		RawResponse:       `{"simulated":true,"status":"SUCCESSFUL"}`,
	}, nil
}

func (s *SimulatedProvider) Status(_ context.Context, providerReference string) (Result, error) {
	return Result{
		Status:            "SUCCESSFUL",
		ProviderStatus:    "SIMULATED_SUCCESS",
		ProviderReference: providerReference,
		RawResponse:       `{"simulated":true,"status":"SUCCESSFUL"}`,
	}, nil
}

func (s *SimulatedProvider) VerifyCallback(_ context.Context, providerReference string) (Result, error) {
	return s.Status(context.Background(), providerReference)
}

func (s *SimulatedProvider) CallbackReference(_ []byte) (string, string, error) {
	return "", "", nil
}
