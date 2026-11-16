package providers

import "testing"

func TestNormalizeMTNStatus(t *testing.T) {
	cases := []struct {
		raw        string
		wantStatus string
		wantReason string
	}{
		{"SUCCESSFUL", "SUCCESSFUL", ""},
		{"FAILED", "FAILED", "FAILED"},
		{"REJECTED", "FAILED", "REJECTED"},
		{"PENDING", "PENDING", ""},
		{"UNKNOWN", "PENDING", ""},
	}
	for _, tc := range cases {
		status, reason := normalizeMTNStatus(tc.raw, "")
		if status != tc.wantStatus || reason != tc.wantReason {
			t.Fatalf("normalizeMTNStatus(%q) = %q/%q, want %q/%q", tc.raw, status, reason, tc.wantStatus, tc.wantReason)
		}
	}
}

func TestNormalizeAirtelStatus(t *testing.T) {
	cases := []struct {
		raw        string
		wantStatus string
	}{
		{"TS", "SUCCESSFUL"},
		{"TF", "FAILED"},
		{"TIP", "PENDING"},
		{"TA", "PENDING"},
		{"", "PENDING"},
	}
	for _, tc := range cases {
		status, _ := normalizeAirtelStatus(tc.raw, "")
		if status != tc.wantStatus {
			t.Fatalf("normalizeAirtelStatus(%q) = %q, want %q", tc.raw, status, tc.wantStatus)
		}
	}
}

func TestAirtelResponseParsing(t *testing.T) {
	body := []byte(`{"data":{"transaction":{"id":"airtel-tx-1","airtel_money_id":"am-1","status":"TS"}},"status":{"code":"200","message":"success"}}`)
	if got := airtelTransactionStatus(body); got != "TS" {
		t.Fatalf("airtelTransactionStatus = %q", got)
	}
	if got := airtelTransactionID(body); got != "am-1" {
		t.Fatalf("airtelTransactionID = %q", got)
	}
}

func TestMTNCallbackReference(t *testing.T) {
	reference, transaction, err := (&MTN{}).CallbackReference([]byte(`{"referenceId":"ref-1","financialTransactionId":"ft-1"}`))
	if err != nil {
		t.Fatalf("callback parse: %v", err)
	}
	if reference != "ref-1" || transaction != "ft-1" {
		t.Fatalf("callback refs = %q/%q", reference, transaction)
	}
}
