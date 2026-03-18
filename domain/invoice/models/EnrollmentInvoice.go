package models

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"

type EnrollmentInvoice struct {
	EnrollmentNumber string            `json:"enrollment_number" gorm:"size:36;not null;primaryKey;index:enrollment_invoice_pk,unique"`
	Enrollment       models.Enrollment `gorm:"foreignKey:EnrollmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	InvoiceNumber    string            `json:"invoice_number" gorm:"size:36;not null;primaryKey;index:enrollment_invoice_pk,unique"`
	Invoice          Invoice           `gorm:"foreignKey:InvoiceNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BalanceDueDate   string            `json:"balance_due_date" gorm:"size:15;"`
}
