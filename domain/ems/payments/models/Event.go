package models

import "time"

const (
	EventPaymentInitiated        = "PAYMENT_INITIATED"
	EventProviderRequestAccepted = "PROVIDER_REQUEST_ACCEPTED"
	EventPaymentPending          = "PAYMENT_PENDING"
	EventCallbackReceived        = "CALLBACK_RECEIVED"
	EventPaymentVerified         = "PAYMENT_VERIFIED"
	EventPaymentSuccessful       = "PAYMENT_SUCCESSFUL"
	EventPaymentFailed           = "PAYMENT_FAILED"
	EventPaymentExpired          = "PAYMENT_EXPIRED"
	EventPaymentCancelled        = "PAYMENT_CANCELLED"
	EventPaymentUnknown          = "PAYMENT_UNKNOWN"
	EventReconciliationStarted   = "RECONCILIATION_STARTED"
	EventReconciliationCompleted = "RECONCILIATION_COMPLETED"
	EventFulfilmentRequested     = "FULFILMENT_REQUESTED"
	EventFulfilmentCompleted     = "FULFILMENT_COMPLETED"
	EventFulfilmentFailed        = "FULFILMENT_FAILED"
)

// PaymentEvent is append-only. It intentionally contains no credentials,
// PINs, access tokens or provider response bodies.
type PaymentEvent struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TransactionNumber string    `json:"transaction_number" gorm:"index;size:36;not null"`
	EventType         string    `json:"event_type" gorm:"index;size:64;not null"`
	ProviderCode      string    `json:"provider_code" gorm:"size:32"`
	ProviderReference string    `json:"provider_reference" gorm:"size:88"`
	ActorUserNumber   string    `json:"actor_user_number" gorm:"size:36"`
	Metadata          string    `json:"metadata" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at"`
}

func (PaymentEvent) TableName() string { return "ems_payment_events" }
