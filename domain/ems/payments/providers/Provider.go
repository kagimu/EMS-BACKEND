package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MTNMoMo     = "MTN_MOMO"
	AirtelMoney = "AIRTEL_MONEY"
)

var (
	ErrNotConfigured = errors.New("mobile-money provider is not configured")
	ErrDisabled      = errors.New("mobile-money provider is disabled")
	ErrContract      = errors.New("mobile-money provider contract is not verified")

	cacheMu sync.Mutex
	cache   = map[string]Provider{}

	testMu        sync.RWMutex
	testProviders = map[string]Provider{}
)

// ProviderError distinguishes definite provider rejections from ambiguous
// network/timeout failures. Ambiguous failures must become UNKNOWN, not FAILED.
type ProviderError struct {
	Operation string
	HTTPCode  int
	Ambiguous bool
	Err       error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Operation + " failed"
}

func (e *ProviderError) Unwrap() error { return e.Err }

func IsAmbiguous(err error) bool {
	var providerErr *ProviderError
	return errors.As(err, &providerErr) && providerErr.Ambiguous
}

type Request struct {
	PaymentNumber string
	Reference     string
	Amount        int64
	Currency      string
	Phone         string
	Description   string
}

type Result struct {
	Status                string
	ProviderStatus        string
	ProviderReference     string
	ProviderTransactionID string
	ProviderAmount        int64
	ProviderCurrency      string
	FailureReason         string
	RawResponse           string
}

type Provider interface {
	Name() string
	DisplayName() string
	Enabled() bool
	Environment() string
	Initiate(ctx context.Context, request Request) (Result, error)
	Status(ctx context.Context, providerReference string) (Result, error)
	VerifyCallback(ctx context.Context, providerReference string) (Result, error)
	CallbackReference(body []byte) (providerReference string, providerTransactionID string, err error)
}

type Info struct {
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	Environment   string `json:"environment"`
	TestMode      bool   `json:"test_mode"`
	ContractReady bool   `json:"contract_ready"`
	Currency      string `json:"currency"`
}

// Get returns the configured provider. Construction is lazy so that .env has
// already been loaded by main. Test providers can be injected with SetForTest.
func Get(name string) (Provider, error) {
	key := strings.ToUpper(strings.TrimSpace(name))
	if key == "" {
		return nil, ErrNotConfigured
	}

	testMu.RLock()
	override, hasOverride := testProviders[key]
	testMu.RUnlock()
	if hasOverride {
		return override, nil
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()
	if existing, ok := cache[key]; ok {
		return existing, nil
	}

	var provider Provider
	switch key {
	case MTNMoMo:
		provider = NewMTNFromEnv()
	case AirtelMoney:
		provider = NewAirtelFromEnv()
	case Simulated:
		provider = NewSimulated()
	default:
		return nil, fmt.Errorf("%w: %s", ErrNotConfigured, name)
	}
	cache[key] = provider
	return provider, nil
}

// EnabledProviders lists providers that are explicitly enabled and have the
// credentials needed to make a request. It never returns secret values.
func EnabledProviders() []Info {
	names := []string{MTNMoMo, AirtelMoney, Simulated}
	result := make([]Info, 0, len(names))
	for _, name := range names {
		provider, err := Get(name)
		if err != nil || !provider.Enabled() {
			continue
		}
		result = append(result, Info{
			Name:          provider.Name(),
			DisplayName:   provider.DisplayName(),
			Environment:   provider.Environment(),
			TestMode:      provider.Environment() != "production",
			ContractReady: true,
			Currency:      "UGX",
		})
	}
	return result
}

// SetForTest installs an in-memory provider override. It is intended for
// integration tests only; production code never calls it.
func SetForTest(name string, provider Provider) func() {
	key := strings.ToUpper(strings.TrimSpace(name))
	testMu.Lock()
	testProviders[key] = provider
	testMu.Unlock()
	return func() {
		testMu.Lock()
		delete(testProviders, key)
		testMu.Unlock()
	}
}

func ResetCacheForTest() {
	cacheMu.Lock()
	cache = map[string]Provider{}
	cacheMu.Unlock()
}

func env(key string) string { return strings.TrimSpace(os.Getenv(key)) }

func envBool(key string, fallback bool) bool {
	value := env(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := env(key)
	if value == "" {
		return fallback
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if parsed, err := time.ParseDuration(value); err == nil && parsed > 0 {
		return parsed
	}
	return fallback
}

func trimURL(value string) string { return strings.TrimRight(strings.TrimSpace(value), "/") }

func firstString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func parseAmount(value interface{}) int64 {
	switch typed := value.(type) {
	case float64:
		return int64(typed)
	case string:
		parsed, _ := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return int64(parsed)
	case json.Number:
		parsed, _ := typed.Int64()
		return parsed
	default:
		return 0
	}
}
