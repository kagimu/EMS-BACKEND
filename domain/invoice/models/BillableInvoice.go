package models

type BillableInvoice struct {
	BillableNumber string   `json:"billable_number" gorm:"size:36;not null;primaryKey;index:billable_invoice_pk,unique"`
	Billable       Billable `gorm:"foreignKey:BillableNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	InvoiceNumber  string   `json:"invoice_number" gorm:"size:36;not null;primaryKey;index:billable_invoice_pk,unique"`
	Invoice        Invoice  `gorm:"foreignKey:InvoiceNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
