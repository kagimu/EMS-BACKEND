package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

// PrepaidBalance holds the current prepaid (wallet) balance for a user (student).
// Money is added via deposits and can later be used to pay invoices.
type PrepaidBalance struct {
	UserNumber string       `json:"user_number" gorm:"size:36;primaryKey;not null"`
	User      authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Balance   float64       `json:"balance" gorm:"not null;default:0"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (PrepaidBalance) TableName() string {
	return "prepaid_balances"
}
