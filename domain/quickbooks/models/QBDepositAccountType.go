package models

import (
	"gorm.io/gorm"
)

type (
	QBDepositAccountType struct {
		gorm.Model
		Number      string `json:"number" gorm:"size:36;unique;not null"`
		Name        string `json:"name" gorm:"size:255;not null"`
		Description string `json:"description" gorm:"type:text"`
		IsActive    bool   `json:"is_active" gorm:"default:true"`
	}
)

func (u *QBDepositAccountType) TableName() string {
	return "qb_deposit_account_types"
}
