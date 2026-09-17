package models

import (
	"gorm.io/gorm"
)

type (
	QBDepositAccount struct {
		gorm.Model
		Number          string         `json:"number" gorm:"size:36;unique;not null"`
		QBProfileNumber string         `json:"qb_profile_number" gorm:"size:36;index"`
		Name            string         `json:"name" gorm:"size:255;not null"`
		AccountTypeID   string         `json:"account_type_id" gorm:"size:36"`
		AccountType     *QBAccountType `json:"account_type,omitempty" gorm:"foreignKey:AccountTypeID;references:Number"`
		Details         string         `json:"details" gorm:"type:text"`
		IsActive        bool           `json:"is_active" gorm:"default:true"`
	}
)

func (u *QBDepositAccount) TableName() string {
	return "qb_deposit_accounts"
}
