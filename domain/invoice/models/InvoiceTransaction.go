package models

type InvoiceTransaction struct {
	InvoiceNumber     string      `json:"invoice_number" gorm:"size:36;not null;primaryKey;index:invoice_transaction_pk,unique"`
	Invoice           Invoice     `gorm:"foreignKey:InvoiceNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	TransactionNumber string      `json:"transaction_number" gorm:"size:36;not null;;primaryKey;index:invoice_transaction_pk,unique"`
	Transaction       Transaction `gorm:"foreignKey:TransactionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AmountPaidOff     float64     `json:"amount_paid_off"`
}
