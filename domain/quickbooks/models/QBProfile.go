package models

import (
	"gorm.io/gorm"
)

type (
	QBProfile struct {
		gorm.Model
		Number           string `json:"number" gorm:"size:36;unique;not null"`
		Name             string `json:"name" gorm:"size:255;not null"`
		Username         string `json:"username" gorm:"size:255;unique;not null"`
		PasswordHash     string `json:"password_hash" gorm:"size:255;not null"`
		DepositToAccount string `json:"deposit_to_account" gorm:"size:255;default:'Undeposited Funds'"`
		IncomeAccount    string `json:"income_account" gorm:"size:255;default:'Service Income'"`
		FileID           string `json:"file_id" gorm:"size:64"`
		OwnerID          string `json:"owner_id" gorm:"size:64"`
		IsDefault        bool   `json:"is_default" gorm:"default:false"`
		IsActive         bool   `json:"is_active" gorm:"default:true"`
		SchoolNumber     string `json:"school_number" gorm:"size:36"`
	}
)

func (u *QBProfile) TableName() string {
	return "qb_profiles"
}
