package models

import (
	"gorm.io/gorm"
)

type (
	QBDepositAccount struct {
		gorm.Model
		Number               string                `json:"number" gorm:"size:36;unique;not null"`
		Name                 string                `json:"name" gorm:"size:255;not null"`
		DepositAccountTypeID string                `json:"deposit_account_type_id" gorm:"size:36"`
		DepositAccountType   *QBDepositAccountType `json:"deposit_account_type,omitempty" gorm:"foreignKey:DepositAccountTypeID;references:Number"`
		Details              string                `json:"details" gorm:"type:text"`
		IsActive             bool                  `json:"is_active" gorm:"default:true"`
	}
)

func (u *QBDepositAccount) TableName() string {
	return "qb_deposit_accounts"
}
