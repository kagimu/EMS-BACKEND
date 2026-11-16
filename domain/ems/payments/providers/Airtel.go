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

const (
	airtelUATBaseURL        = "https://openapiuat.airtel.africa"
	airtelProductionBaseURL = "https://openapi.airtel.africa"
)

type AirtelConfig struct {
	Enabled          bool
	BaseURL          string
	ClientID         string
	ClientSecret     string
	Country          string
	Currency         string
	CallbackURL      string
	ContractVerified bool
	Timeout          time.Duration
}

func AirtelConfigFromEnv() AirtelConfig {
	environment := strings.ToLower(env("AIRTEL_MONEY_ENVIRONMENT"))
	if environment == "" {
		environment = "uat"
	}
	baseURL := trimURL(env("AIRTEL_MONEY_BASE_URL"))
	if baseURL == "" {
		if environment == "production" {
			baseURL = airtelProductionBaseURL
		} else {
			baseURL = airtelUATBaseURL
		}
	}
	country := strings.ToUpper(env("AIRTEL_MONEY_COUNTRY"))
	if country == "" {
		country = "UG"
	}
	currency := strings.ToUpper(env("AIRTEL_MONEY_CURRENCY"))
	if currency == "" {
		currency = "UGX"
	}
	return AirtelConfig{
		Enabled:          envBool("AIRTEL_MONEY_ENABLED", false),
		BaseURL:          baseURL,
		ClientID:         env("AIRTEL_MONEY_CLIENT_ID"),
		ClientSecret:     env("AIRTEL_MONEY_CLIENT_SECRET"),
		Country:          country,
		Currency:         currency,
		CallbackURL:      env("AIRTEL_MONEY_CALLBACK_URL"),
		ContractVerified: envBool("AIRTEL_MONEY_CONTRACT_VERIFIED", false),
		Timeout:          envDuration("AIRTEL_MONEY_TIMEOUT_SECONDS", 25*time.Second),
	}
}

type Airtel struct {
	cfg    AirtelConfig
	client *http.Client

	tokenMu     sync.Mutex
	token       string
	tokenExpiry time.Time
}

func NewAirtelFromEnv() *Airtel {
	cfg := AirtelConfigFromEnv()
	return &Airtel{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

func (a *Airtel) Name() string { return AirtelMoney }

func (a *Airtel) DisplayName() string { return "Airtel Money" }

func (a *Airtel) Environment() string {
	if strings.Contains(strings.ToLower(a.cfg.BaseURL), "uat") {
		return "uat"
	}
	return "production"
}

func (a *Airtel) Enabled() bool {
	return a.cfg.Enabled && a.cfg.ContractVerified && a.cfg.BaseURL != "" && a.cfg.ClientID != "" && a.cfg.ClientSecret != ""
}

func (a *Airtel) accessToken(ctx context.Context) (string, error) {
	if !a.Enabled() {
		if a.cfg.Enabled && !a.cfg.ContractVerified {
			return "", fmt.Errorf("%w: Airtel Uganda API contract has not been verified", ErrContract)
		}
		if a.cfg.Enabled {
			return "", fmt.Errorf("%w: Airtel Money credentials are incomplete", ErrNotConfigured)
		}
		return "", ErrDisabled
	}
	a.tokenMu.Lock()
	defer a.tokenMu.Unlock()
	if a.token != "" && time.Now().Before(a.tokenExpiry) {
		return a.token, nil
	}

	payload, _ := json.Marshal(map[string]string{
		"client_id":     a.cfg.ClientID,
		"client_secret": a.cfg.ClientSecret,
		"grant_type":    "client_credentials",
	})
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, a.cfg.BaseURL+"/auth/oauth2/token", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	response, err := a.client.Do(request)
	if err != nil {
		return "", &ProviderError{Operation: "Airtel token", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", &ProviderError{Operation: "Airtel token", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("Airtel token returned HTTP %d", response.StatusCode)}
	}
	var parsed struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("Airtel Money token response was invalid: %w", err)
	}
	if parsed.AccessToken == "" {
		return "", errors.New("Airtel Money token response did not include an access token")
	}
	lifetime := time.Duration(parsed.ExpiresIn) * time.Second
	if lifetime <= 0 {
		lifetime = 50 * time.Minute
	}
	if lifetime > time.Minute {
		lifetime -= time.Minute
	}
	a.token = parsed.AccessToken
	a.tokenExpiry = time.Now().Add(lifetime)
	return a.token, nil
}

func (a *Airtel) Initiate(ctx context.Context, request Request) (Result, error) {
	token, err := a.accessToken(ctx)
	if err != nil {
		return Result{}, err
	}
	payload := map[string]interface{}{
		"reference": request.Reference,
		"subscriber": map[string]string{
			"country":  a.cfg.Country,
			"currency": a.cfg.Currency,
			"msisdn":   request.Phone,
		},
		"transaction": map[string]interface{}{
			"amount":   request.Amount,
			"country":  a.cfg.Country,
			"currency": a.cfg.Currency,
			"id":       request.PaymentNumber,
		},
		"additional_info": map[string]string{
			"message": truncate(request.Description, 160),
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.cfg.BaseURL+"/merchant/v1/payments/", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Country", a.cfg.Country)
	req.Header.Set("X-Currency", a.cfg.Currency)

	response, err := a.client.Do(req)
	if err != nil {
		return Result{}, &ProviderError{Operation: "Airtel collection", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Result{}, &ProviderError{Operation: "Airtel collection", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 408 || response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("Airtel collection returned HTTP %d: %s", response.StatusCode, truncate(string(raw), 300))}
	}
	providerStatus := airtelTransactionStatus(raw)
	status, failure := normalizeAirtelStatus(providerStatus, airtelFailureReason(raw))
	return Result{
		Status:                status,
		ProviderStatus:        providerStatus,
		ProviderReference:     request.PaymentNumber,
		ProviderTransactionID: airtelTransactionID(raw),
		ProviderAmount:        airtelTransactionAmount(raw),
		ProviderCurrency:      airtelTransactionCurrency(raw),
		FailureReason:         failure,
		RawResponse:           string(raw),
	}, nil
}

func (a *Airtel) Status(ctx context.Context, providerReference string) (Result, error) {
	token, err := a.accessToken(ctx)
	if err != nil {
		return Result{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.cfg.BaseURL+"/standard/v1/payments/"+providerReference, nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Country", a.cfg.Country)
	req.Header.Set("X-Currency", a.cfg.Currency)

	response, err := a.client.Do(req)
	if err != nil {
		return Result{}, &ProviderError{Operation: "Airtel status", Ambiguous: true, Err: err}
	}
	defer response.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Result{}, &ProviderError{Operation: "Airtel status", HTTPCode: response.StatusCode, Ambiguous: response.StatusCode == 408 || response.StatusCode == 429 || response.StatusCode >= 500, Err: fmt.Errorf("Airtel status returned HTTP %d: %s", response.StatusCode, truncate(string(raw), 300))}
	}
	providerStatus := airtelTransactionStatus(raw)
	status, failure := normalizeAirtelStatus(providerStatus, airtelFailureReason(raw))
	return Result{
		Status:                status,
		ProviderStatus:        providerStatus,
		ProviderReference:     providerReference,
		ProviderTransactionID: airtelTransactionID(raw),
		ProviderAmount:        airtelTransactionAmount(raw),
		ProviderCurrency:      airtelTransactionCurrency(raw),
		FailureReason:         failure,
		RawResponse:           string(raw),
	}, nil
}

// VerifyCallback performs an independent status lookup. The Airtel callback
// contract is not verified, so the adapter remains disabled until explicitly
// configured with a verified contract flag.
func (a *Airtel) VerifyCallback(ctx context.Context, providerReference string) (Result, error) {
	return a.Status(ctx, providerReference)
}

func (a *Airtel) CallbackReference(body []byte) (string, string, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", fmt.Errorf("invalid Airtel Money callback: %w", err)
	}
	reference := firstString(
		stringValue(payload["reference"]),
		stringValue(payload["client_transaction_id"]),
		nestedString(payload, "transaction", "id"),
		nestedString(payload, "data", "transaction", "id"),
	)
	transactionID := firstString(
		stringValue(payload["airtel_money_id"]),
		nestedString(payload, "transaction", "airtel_money_id"),
		nestedString(payload, "data", "transaction", "airtel_money_id"),
	)
	return reference, transactionID, nil
}

func airtelTransactionAmount(raw []byte) int64 {
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0
	}
	return parseAmount(firstValue(
		nestedValue(payload, "data", "transaction", "amount"),
		nestedValue(payload, "transaction", "amount"),
	))
}

func airtelTransactionCurrency(raw []byte) string {
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return firstString(
		nestedString(payload, "data", "transaction", "currency"),
		nestedString(payload, "transaction", "currency"),
	)
}

func nestedValue(payload map[string]interface{}, path ...string) interface{} {
	var current interface{} = payload
	for _, key := range path {
		object, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current = object[key]
	}
	return current
}

func firstValue(values ...interface{}) interface{} {
	for _, value := range values {
		if value != nil && stringValue(value) != "" {
			return value
		}
	}
	return nil
}

func normalizeAirtelStatus(raw, reason string) (string, string) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "TS", "SUCCESS", "SUCCESSFUL", "SUCCESSFULL", "COMPLETED":
		return "SUCCESSFUL", ""
	case "TF", "FAILED", "FAILURE", "DECLINED", "CANCELLED", "REVERSED":
		return "FAILED", firstString(reason, raw)
	case "TA", "TIP", "PENDING", "PROCESSING", "IN_PROGRESS", "INITIATED":
		return "PENDING", ""
	default:
		return "PENDING", ""
	}
}

func airtelTransactionStatus(raw []byte) string {
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return firstString(
		nestedString(payload, "data", "transaction", "status"),
		nestedString(payload, "transaction", "status"),
		stringValue(payload["status"]),
	)
}

func airtelFailureReason(raw []byte) string {
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return firstString(
		nestedString(payload, "status", "message"),
		nestedString(payload, "data", "transaction", "message"),
	)
}

func airtelTransactionID(raw []byte) string {
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return firstString(
		nestedString(payload, "data", "transaction", "airtel_money_id"),
		nestedString(payload, "transaction", "airtel_money_id"),
		nestedString(payload, "data", "transaction", "id"),
	)
}

func nestedString(payload map[string]interface{}, path ...string) string {
	var current interface{} = payload
	for _, key := range path {
		object, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		current = object[key]
	}
	return stringValue(current)
}

func stringValue(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return fmt.Sprintf("%.0f", typed)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}
