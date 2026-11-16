package models

import "time"

const (
	ReconciliationPending             = "PENDING"
	ReconciliationMatched             = "MATCHED"
	ReconciliationMismatch            = "MISMATCH"
	ReconciliationProviderUnavailable = "PROVIDER_UNAVAILABLE"
	ReconciliationUnknown             = "UNKNOWN"
)

type PaymentReconciliation struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TransactionNumber string    `json:"transaction_number" gorm:"index;size:36;not null"`
	ProviderCode      string    `json:"provider_code" gorm:"index;size:32;not null"`
	ProviderReference string    `json:"provider_reference" gorm:"size:88"`
	ProviderAmount    int64     `json:"provider_amount"`
	EMSAmount         int64     `json:"ems_amount"`
	ProviderCurrency  string    `json:"provider_currency" gorm:"size:8"`
	EMSCurrency       string    `json:"ems_currency" gorm:"size:8"`
	Result            string    `json:"result" gorm:"index;size:32"`
	Discrepancy       string    `json:"discrepancy" gorm:"type:text"`
	AttemptNumber     int       `json:"attempt_number"`
	ReconciledAt      time.Time `json:"reconciled_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (PaymentReconciliation) TableName() string { return "ems_payment_reconciliations" }
