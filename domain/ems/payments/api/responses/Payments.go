package responses

import "time"

type PaymentView struct {
	Number                string     `json:"number"`
	Provider              string     `json:"provider"`
	ProviderName          string     `json:"provider_name"`
	Purpose               string     `json:"purpose"`
	Reference             string     `json:"reference"`
	OrganisationNumber    string     `json:"organisation_number"`
	OrganisationName      string     `json:"organisation_name"`
	CommercialOrderNumber string     `json:"commercial_order_number"`
	CustomerProfileNumber string     `json:"customer_profile_number"`
	CustomerName          string     `json:"customer_name,omitempty"`
	CustomerEmail         string     `json:"customer_email,omitempty"`
	PayerPhoneMasked      string     `json:"payer_phone_masked"`
	AmountUGX             int64      `json:"amount_ugx"`
	Currency              string     `json:"currency"`
	Status                string     `json:"status"`
	ProviderStatus        string     `json:"provider_status"`
	ProviderReference     string     `json:"provider_reference"`
	ProviderTransactionID string     `json:"provider_transaction_id"`
	FailureCode           string     `json:"failure_code"`
	FailureReason         string     `json:"failure_reason"`
	FulfilmentStatus      string     `json:"fulfilment_status"`
	FulfilmentError       string     `json:"fulfilment_error,omitempty"`
	ExpiresAt             *time.Time `json:"expires_at"`
	UnknownAt             *time.Time `json:"unknown_at"`
	InitiatedAt           time.Time  `json:"initiated_at"`
	LastCheckedAt         *time.Time `json:"last_checked_at"`
	PaidAt                *time.Time `json:"paid_at"`
	VerifiedAt            *time.Time `json:"verified_at"`
	TokenClaimedAt        *time.Time `json:"token_claimed_at"`
	TokenClaimed          bool       `json:"token_claimed"`
}

type PaymentEventView struct {
	EventType         string    `json:"event_type"`
	ProviderCode      string    `json:"provider_code"`
	ProviderReference string    `json:"provider_reference"`
	ActorUserNumber   string    `json:"actor_user_number"`
	Metadata          string    `json:"metadata"`
	CreatedAt         time.Time `json:"created_at"`
}

type ReconciliationView struct {
	Result            string    `json:"result"`
	ProviderReference string    `json:"provider_reference"`
	ProviderAmount    int64     `json:"provider_amount"`
	EMSAmount         int64     `json:"ems_amount"`
	ProviderCurrency  string    `json:"provider_currency"`
	EMSCurrency       string    `json:"ems_currency"`
	Discrepancy       string    `json:"discrepancy"`
	AttemptNumber     int       `json:"attempt_number"`
	ReconciledAt      time.Time `json:"reconciled_at"`
}

type PaymentDetail struct {
	PaymentView
	PayerPhone       string               `json:"payer_phone"`
	CheckCount       int                  `json:"check_count"`
	VerifiedByNumber string               `json:"verified_by_number"`
	OrderStatus      string               `json:"order_status"`
	Events           []PaymentEventView   `json:"events"`
	Reconciliations  []ReconciliationView `json:"reconciliations"`
}

type ProviderTotal struct {
	Provider  string `json:"provider"`
	Count     int64  `json:"count"`
	AmountUGX int64  `json:"amount_ugx"`
}

type PaymentSummary struct {
	TotalCount          int64           `json:"total_count"`
	SuccessfulCount     int64           `json:"successful_count"`
	PendingCount        int64           `json:"pending_count"`
	UnknownCount        int64           `json:"unknown_count"`
	FailedCount         int64           `json:"failed_count"`
	SuccessfulAmountUGX int64           `json:"successful_amount_ugx"`
	ByProvider          []ProviderTotal `json:"by_provider"`
}

type PaymentList struct {
	Items   []PaymentView            `json:"items"`
	Total   int64                    `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Summary PaymentSummary           `json:"summary"`
	Worker  ReconciliationWorkerView `json:"worker"`
}

type ReconciliationWorkerView struct {
	Enabled          bool      `json:"enabled"`
	Running          bool      `json:"running"`
	IntervalSeconds  int       `json:"interval_seconds"`
	MaxAttempts      int       `json:"max_attempts"`
	LastRunAt        time.Time `json:"last_run_at"`
	LastSuccessAt    time.Time `json:"last_success_at"`
	LastError        string    `json:"last_error"`
	LastScannedCount int       `json:"last_scanned_count"`
	LastUpdatedCount int       `json:"last_updated_count"`
}

type ProviderInfo struct {
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	Environment   string `json:"environment"`
	TestMode      bool   `json:"test_mode"`
	ContractReady bool   `json:"contract_ready"`
	Currency      string `json:"currency"`
}

type PaymentOptions struct {
	Providers         []ProviderInfo `json:"providers"`
	SimulatedEnabled  bool           `json:"simulated_enabled"`
	Currency          string         `json:"currency"`
	MaxAmountUGX      int64          `json:"max_amount_ugx"`
	HasActiveProvider bool           `json:"has_active_provider"`
}
