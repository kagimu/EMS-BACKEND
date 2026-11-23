package models

import (
	"gorm.io/gorm"
)

type (
	QBSettings struct {
		gorm.Model
		Number           string `json:"number" gorm:"size:36;unique;not null"`
		SchoolNumber     string `json:"school_number" gorm:"size:36;not null;unique"`
		Username         string `json:"username" gorm:"size:255;not null"`
		PasswordHash     string `json:"password_hash" gorm:"size:255;not null"`
		IsActive         bool   `json:"is_active" gorm:"default:true"`
		DepositToAccount string `json:"deposit_to_account" gorm:"size:255;default:'Undeposited Funds'"`
		IncomeAccount    string `json:"income_account" gorm:"size:255;default:'Service Income'"`
	}
)

func (u *QBSettings) TableName() string {
	return "qb_settings"
}
