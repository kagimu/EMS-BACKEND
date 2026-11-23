package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const mtnSandboxBaseURL = "https://sandbox.momodeveloper.mtn.com"

type MTNConfig struct {
	Enabled         bool
	BaseURL         string
	SubscriptionKey string
	APIUser         string
	APIKey          string
	TargetEnv       string
	Currency        string
	CallbackURL     string
	Timeout         time.Duration
}

func MTNConfigFromEnv() MTNConfig {
	environment := strings.ToLower(env("MTN_MOMO_ENVIRONMENT"))
	if environment == "" {
		environment = "sandbox"
	}
	baseURL := trimURL(env("MTN_MOMO_BASE_URL"))
	if baseURL == "" && environment == "sandbox" {
		baseURL = mtnSandboxBaseURL
	}
	target := env("MTN_MOMO_TARGET_ENVIRONMENT")
	if target == "" && environment == "sandbox" {
		target = "sandbox"
	}
	currency := env("MTN_MOMO_CURRENCY")
	if currency == "" {
		currency = "UGX"
	}
	return MTNConfig{
		Enabled:         envBool("MTN_MOMO_ENABLED", false),
		BaseURL:         baseURL,
		SubscriptionKey: env("MTN_MOMO_SUBSCRIPTION_KEY"),
		APIUser:         env("MTN_MOMO_API_USER"),
		APIKey:          env("MTN_MOMO_API_KEY"),
		TargetEnv:       target,
		Currency:        currency,
		CallbackURL:     env("MTN_MOMO_CALLBACK_URL"),
		Timeout:         envDuration("MTN_MOMO_TIMEOUT_SECONDS", 25*time.Second),
	}
}

type MTN struct {
	cfg    MTNConfig
	client *http.Client

	tokenMu     sync.Mutex
	token       string
	tokenExpiry time.Time
}

func NewMTNFromEnv() *MTN {
	cfg := MTNConfigFromEnv()
	return &MTN{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

func (m *MTN) Name() string { return MTNMoMo }

func (m *MTN) DisplayName() string { return "MTN MoMo" }

func (m *MTN) Environment() string {
	if m.cfg.TargetEnv == "" || strings.EqualFold(m.cfg.TargetEnv, "sandbox") {
		return "sandbox"
	}
	return "production"
}

func (m *MTN) Enabled() bool {
	return m.cfg.Enabled && m.cfg.BaseURL != "" && m.cfg.SubscriptionKey != "" &&
		m.cfg.APIUser != "" && m.cfg.APIKey != "" && m.cfg.TargetEnv != ""
}

func (m *MTN) accessToken(ctx context.Context) (string, error) {
	if !m.Enabled() {
		if m.cfg.Enabled {
			return "", fmt.Errorf("%w: MTN MoMo credentials are incomplete", ErrNotConfigured)
		}
		return "", ErrDisabled
	}
	m.tokenMu.Lock()
	defer m.tokenMu.Unlock()
	if m.token != "" && time.Now().Before(m.tokenExpiry) {
		return m.token, nil
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, m.cfg.BaseURL+"/collection/token/", nil)
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(m.cfg.APIUser, m.cfg.APIKey)
	request.Header.Set("Ocp-Apim-Subscription-Key", m.cfg.SubscriptionKey)
	request.Header.Set("X-Target-Environment", m.cfg.TargetEnv)
	request.Header.Set("Accept", "application/json")

	response, err := m.client.Do(request)
	if err != nil {
		return "", &ProviderError{Operation: "MTN token", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", &ProviderError{Operation: "MTN token", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("MTN token returned HTTP %d", response.StatusCode)}
	}
	var parsed struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("MTN MoMo token response was invalid: %w", err)
	}
	if parsed.AccessToken == "" {
		return "", errors.New("MTN MoMo token response did not include an access token")
	}
	lifetime := time.Duration(parsed.ExpiresIn) * time.Second
	if lifetime <= 0 {
		lifetime = 50 * time.Minute
	}
	if lifetime > time.Minute {
		lifetime -= time.Minute
	}
	m.token = parsed.AccessToken
	m.tokenExpiry = time.Now().Add(lifetime)
	return m.token, nil
}

func (m *MTN) Initiate(ctx context.Context, request Request) (Result, error) {
	token, err := m.accessToken(ctx)
	if err != nil {
		return Result{}, err
	}
	payload := map[string]interface{}{
		"amount":     fmt.Sprintf("%d", request.Amount),
		"currency":   request.Currency,
		"externalId": request.Reference,
		"payer": map[string]string{
			"partyIdType": "MSISDN",
			"partyId":     request.Phone,
		},
		"payerMessage": truncate(request.Description, 160),
		"payeeNote":    "EMS software payment",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.cfg.BaseURL+"/collection/v1_0/requesttopay", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.cfg.SubscriptionKey)
	req.Header.Set("X-Reference-Id", request.PaymentNumber)
	req.Header.Set("X-Target-Environment", m.cfg.TargetEnv)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if m.cfg.CallbackURL != "" {
		req.Header.Set("X-Callback-Url", m.cfg.CallbackURL)
	}

	response, err := m.client.Do(req)
	if err != nil {
		return Result{}, &ProviderError{Operation: "MTN requesttopay", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusOK {
		return Result{}, &ProviderError{Operation: "MTN requesttopay", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 408 || response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("MTN requesttopay returned HTTP %d: %s", response.StatusCode, truncate(string(raw), 300))}
	}
	return Result{
		Status:            "PENDING",
		ProviderStatus:    "PENDING",
		ProviderReference: request.PaymentNumber,
		RawResponse:       string(raw),
	}, nil
}

func (m *MTN) Status(ctx context.Context, providerReference string) (Result, error) {
	token, err := m.accessToken(ctx)
	if err != nil {
		return Result{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.cfg.BaseURL+"/collection/v1_0/requesttopay/"+providerReference, nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.cfg.SubscriptionKey)
	req.Header.Set("X-Target-Environment", m.cfg.TargetEnv)
	req.Header.Set("Accept", "application/json")

	response, err := m.client.Do(req)
	if err != nil {
		return Result{}, &ProviderError{Operation: "MTN status", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Result{}, &ProviderError{Operation: "MTN status", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 408 || response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("MTN status returned HTTP %d: %s", response.StatusCode, truncate(string(raw), 300))}
	}
	var parsed struct {
		Status                 string      `json:"status"`
		Reason                 string      `json:"reason"`
		FinancialTransactionID string      `json:"financialTransactionId"`
		Amount                 interface{} `json:"amount"`
		Currency               string      `json:"currency"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Result{}, fmt.Errorf("MTN MoMo status response was invalid: %w", err)
	}
	status, failure := normalizeMTNStatus(parsed.Status, parsed.Reason)
	return Result{
		Status:                status,
		ProviderStatus:        parsed.Status,
		ProviderReference:     providerReference,
		ProviderTransactionID: parsed.FinancialTransactionID,
		ProviderAmount:        parseAmount(parsed.Amount),
		ProviderCurrency:      parsed.Currency,
		FailureReason:         failure,
		RawResponse:           string(raw),
	}, nil
}

// VerifyCallback performs an independent status lookup. MTN does not document
// a callback signature; the callback body is never authoritative on its own.
func (m *MTN) VerifyCallback(ctx context.Context, providerReference string) (Result, error) {
	return m.Status(ctx, providerReference)
}

func (m *MTN) CallbackReference(body []byte) (string, string, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", fmt.Errorf("invalid MTN MoMo callback: %w", err)
	}
	reference := firstString(stringValue(payload["referenceId"]), stringValue(payload["reference_id"]))
	transactionID := firstString(stringValue(payload["financialTransactionId"]), stringValue(payload["financial_transaction_id"]))
	return reference, transactionID, nil
}

func normalizeMTNStatus(raw, reason string) (string, string) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "SUCCESSFUL", "SUCCESS":
		return "SUCCESSFUL", ""
	case "FAILED", "REJECTED", "TIMEOUT", "CANCELLED":
		return "FAILED", firstString(reason, raw)
	case "PENDING", "PROCESSING", "ONGOING":
		return "PENDING", ""
	default:
		return "PENDING", ""
	}
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
