package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

type BillPaymentPlan struct {
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
	Number               string                       `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	BillableNumber       *string                      `json:"billable_number" gorm:"size:36;index;default:null"`
	PromissoryNoteNumber *string                      `json:"promissory_note_number" gorm:"size:36;index;default:null"`
	PromissoryNote       Invoice                      `gorm:"foreignKey:PromissoryNoteNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber        string                       `json:"student_number" gorm:"size:36;not null;index"`
	Student              studentModel.Student         `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber    string                       `json:"study_period_number" gorm:"size:36;not null;index"`
	TotalAmount          float64                      `json:"total_amount"`
	StartDate            time.Time                    `json:"start_date" gorm:"type:date"`
	Notes                string                       `json:"notes" gorm:"type:text"`
	Status               string                       `json:"status" gorm:"size:36;not null;default:ACTIVE;index"`
	UserNumber           string                       `json:"user_number" gorm:"size:36;not null"`
	User                 authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Installments         []BillPaymentPlanInstallment `json:"installments" gorm:"foreignKey:BillPaymentPlanNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type BillPaymentPlanInstallment struct {
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Number                string    `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	BillPaymentPlanNumber string    `json:"bill_payment_plan_number" gorm:"size:36;not null;index"`
	DueDate               time.Time `json:"due_date" gorm:"type:date"`
	Amount                float64   `json:"amount"`
	PaidAmount            float64   `json:"paid_amount"`
	Status                string    `json:"status" gorm:"size:36;not null;default:PENDING"`
	Sequence              int       `json:"sequence"`
}
