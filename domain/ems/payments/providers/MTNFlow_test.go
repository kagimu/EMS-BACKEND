package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMTNOAuthAndRequestToPayFlowWithMockServer(t *testing.T) {
	var tokenCalled, requestCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/collection/token/":
			tokenCalled = true
			if r.Method != http.MethodPost {
				t.Errorf("token method = %s", r.Method)
			}
			user, password, ok := r.BasicAuth()
			if !ok || user != "api-user" || password != "api-key" {
				t.Errorf("unexpected basic auth: %s/%s ok=%v", user, password, ok)
			}
			if r.Header.Get("Ocp-Apim-Subscription-Key") != "subscription-key" {
				t.Errorf("missing subscription key")
			}
			if r.Header.Get("X-Target-Environment") != "sandbox" {
				t.Errorf("missing target environment")
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"token-123","expires_in":3600}`))
		case "/collection/v1_0/requesttopay":
			requestCalled = true
			if r.Method != http.MethodPost {
				t.Errorf("requesttopay method = %s", r.Method)
			}
			if r.Header.Get("Authorization") != "Bearer token-123" {
				t.Errorf("missing bearer token")
			}
			if r.Header.Get("X-Reference-Id") != "payment-reference" {
				t.Errorf("missing X-Reference-Id")
			}
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["currency"] != "UGX" || body["amount"] != "5000" {
				t.Errorf("unexpected request body: %#v", body)
			}
			w.WriteHeader(http.StatusAccepted)
		case "/collection/v1_0/requesttopay/payment-reference":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"SUCCESSFUL","amount":"5000","currency":"UGX","financialTransactionId":"ft-1"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := &MTN{
		cfg: MTNConfig{
			Enabled: true, BaseURL: server.URL, SubscriptionKey: "subscription-key",
			APIUser: "api-user", APIKey: "api-key", TargetEnv: "sandbox",
			Currency: "UGX", Timeout: 5 * time.Second,
		},
		client: server.Client(),
	}
	result, err := provider.Initiate(context.Background(), Request{
		PaymentNumber: "payment-reference", Reference: "order-1",
		Amount: 5000, Currency: "UGX", Phone: "256771234567", Description: "test",
	})
	if err != nil {
		t.Fatalf("initiate: %v", err)
	}
	if result.Status != "PENDING" || !tokenCalled || !requestCalled {
		t.Fatalf("unexpected initiate result or missing calls: %#v token=%v request=%v", result, tokenCalled, requestCalled)
	}
	status, err := provider.Status(context.Background(), "payment-reference")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if status.Status != "SUCCESSFUL" || status.ProviderAmount != 5000 || status.ProviderCurrency != "UGX" {
		t.Fatalf("unexpected status: %#v", status)
	}
}
