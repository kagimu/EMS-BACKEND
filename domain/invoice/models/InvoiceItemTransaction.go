package models

import (
	"time"
)

type InvoiceItemTransaction struct {
	Number            string      `json:"number" gorm:"size:36;primaryKey;unique;not null"`
	TransactionNumber string      `json:"transaction_number" gorm:"size:36;not null;index"`
	Transaction       Transaction `gorm:"foreignKey:TransactionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	InvoiceItemNumber string      `json:"invoice_item_number" gorm:"size:36;not null;index"`
	InvoiceItem       InvoiceItem `gorm:"foreignKey:InvoiceItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AmountPaidOff     float64     `json:"amount_paid_off" gorm:"not null"`
	CreatedAt         time.Time   `json:"createdAt" gorm:"not null"`
	ItemReceiptNumber string      `json:"item_receipt_no" gorm:"size:36;unique"`
}

func (u *InvoiceItemTransaction) TableName() string {
	return "invoice_item_transactions"
}
