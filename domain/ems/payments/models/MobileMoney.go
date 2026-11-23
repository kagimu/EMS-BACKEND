package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ProviderMTNMoMo     = "MTN_MOMO"
	ProviderAirtelMoney = "AIRTEL_MONEY"
	ProviderSimulated   = "SIMULATED"

	PurposeSetup        = "SETUP"
	PurposeSubscription = "SUBSCRIPTION"

	StatusInitiated  = "INITIATED"
	StatusPending    = "PENDING"
	StatusSuccessful = "SUCCESSFUL"
	StatusFailed     = "FAILED"
	StatusExpired    = "EXPIRED"
	StatusCancelled  = "CANCELLED"
	StatusUnknown    = "UNKNOWN"

	FulfilmentNotRequired = "NOT_REQUIRED"
	FulfilmentPending     = "PENDING"
	FulfilmentCompleted   = "COMPLETED"
	FulfilmentFailed      = "FAILED"
)

var PaymentStatuses = []string{
	StatusInitiated, StatusPending, StatusSuccessful, StatusFailed,
	StatusExpired, StatusCancelled, StatusUnknown,
}

// MobileMoneyPayment is the authoritative record of a mobile-money collection
// attempt. Amounts are copied from the server-side commercial order. Provider
// credentials and access tokens are never stored here.
type MobileMoneyPayment struct {
	gorm.Model
	Number                string     `json:"number" gorm:"uniqueIndex;size:36"`
	Provider              string     `json:"provider" gorm:"index;size:32"`
	Purpose               string     `json:"purpose" gorm:"index;size:24"`
	Reference             string     `json:"reference" gorm:"index;size:36"`
	CommercialOrderNumber string     `json:"commercial_order_number" gorm:"index;size:36"`
	OrganisationNumber    string     `json:"organisation_number" gorm:"index;size:36"`
	CustomerProfileNumber string     `json:"customer_profile_number" gorm:"index;size:36"`
	PayerPhone            string     `json:"payer_phone" gorm:"size:32"`
	AmountUGX             int64      `json:"amount_ugx"`
	Currency              string     `json:"currency" gorm:"size:8;default:UGX"`
	Status                string     `json:"status" gorm:"index;size:24"`
	StatusVersion         int        `json:"status_version" gorm:"not null;default:1"`
	ProviderReference     string     `json:"provider_reference" gorm:"uniqueIndex;size:88"`
	ProviderTransactionID string     `json:"provider_transaction_id" gorm:"index;size:128"`
	ProviderStatus        string     `json:"provider_status" gorm:"size:48"`
	ProviderAmount        int64      `json:"provider_amount"`
	ProviderCurrency      string     `json:"provider_currency" gorm:"size:8"`
	MerchantReference     *string    `json:"merchant_reference,omitempty" gorm:"uniqueIndex;size:88"`
	IdempotencyKey        *string    `json:"idempotency_key,omitempty" gorm:"uniqueIndex;size:191"`
	FailureCode           string     `json:"failure_code" gorm:"size:64"`
	FailureReason         string     `json:"failure_reason" gorm:"type:text"`
	RequestJSON           string     `json:"request_json" gorm:"type:text"`
	ResponseJSON          string     `json:"response_json" gorm:"type:text"`
	CallbackJSON          string     `json:"callback_json" gorm:"type:text"`
	InitiatedByNumber     string     `json:"initiated_by_number" gorm:"size:36"`
	CreatedByNumber       string     `json:"created_by_number" gorm:"size:36"`
	ExpiresAt             *time.Time `json:"expires_at"`
	LastCheckedAt         *time.Time `json:"last_checked_at"`
	CheckCount            int        `json:"check_count"`
	PaidAt                *time.Time `json:"paid_at"`
	UnknownAt             *time.Time `json:"unknown_at"`
	VerifiedAt            *time.Time `json:"verified_at"`
	VerifiedByNumber      string     `json:"verified_by_number" gorm:"size:36"`
	FulfilmentStatus      string     `json:"fulfilment_status" gorm:"index;size:24"`
	FulfilmentError       string     `json:"fulfilment_error" gorm:"type:text"`
	FulfilmentCompletedAt *time.Time `json:"fulfilment_completed_at"`
	TokenClaimedAt        *time.Time `json:"token_claimed_at"`
	ActivationTokenNumber string     `json:"activation_token_number" gorm:"size:36"`
}

func (MobileMoneyPayment) TableName() string { return "ems_mobile_money_payments" }

func IsFinal(status string) bool {
	switch status {
	case StatusSuccessful, StatusFailed, StatusCancelled, StatusExpired:
		return true
	default:
		return false
	}
}

func IsUnresolved(status string) bool {
	return status == StatusInitiated || status == StatusPending || status == StatusUnknown
}

// CanTransition is the central status transition guard. SUCCESSFUL is
// terminal; UNKNOWN is not final but must be resolved before a new charge.
func CanTransition(from, to string) bool {
	from = normalizeStatus(from)
	to = normalizeStatus(to)
	if from == "" {
		return to != ""
	}
	if from == to {
		return true
	}
	if from == StatusSuccessful {
		return false
	}
	allowed := map[string]map[string]bool{
		StatusInitiated: {
			StatusPending: true, StatusUnknown: true, StatusFailed: true,
			StatusCancelled: true, StatusExpired: true,
		},
		StatusPending: {
			StatusUnknown: true, StatusSuccessful: true, StatusFailed: true,
			StatusCancelled: true, StatusExpired: true,
		},
		StatusUnknown: {
			StatusPending: true, StatusSuccessful: true, StatusFailed: true,
			StatusCancelled: true, StatusExpired: true,
		},
	}
	return allowed[from][to]
}

func normalizeStatus(value string) string {
	switch value {
	case StatusInitiated, StatusPending, StatusSuccessful, StatusFailed, StatusExpired, StatusCancelled, StatusUnknown:
		return value
	default:
		return ""
	}
}
