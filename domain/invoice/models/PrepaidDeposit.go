package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

// PrepaidDeposit records a single deposit (money brought in) for audit trail.
type PrepaidDeposit struct {
	Number          string       `json:"number" gorm:"size:36;primaryKey;not null"`
	UserNumber      string       `json:"user_number" gorm:"size:36;not null;index"`
	User            authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount          float64      `json:"amount" gorm:"not null"`
	DatePaid        time.Time    `json:"date_paid" gorm:"not null"`
	Description     string       `json:"description" gorm:"size:500"`
	ClientReference string       `json:"client_reference" gorm:"size:62"`
	AddedByNumber   string       `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy         authModel.User `gorm:"foreignKey:AddedByNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt       time.Time    `json:"created_at"`
}

func (PrepaidDeposit) TableName() string {
	return "prepaid_deposits"
}
