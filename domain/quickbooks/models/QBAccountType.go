package models

import (
	"gorm.io/gorm"
)

type (
	QBAccountType struct {
		gorm.Model
		Number       string `json:"number" gorm:"size:36;unique;not null"`
		Name         string `json:"name" gorm:"size:255;not null"`
		Description  string `json:"description" gorm:"type:text"`
		IsActive     bool   `json:"is_active" gorm:"default:true"`
	}
)

func (u *QBAccountType) TableName() string {
	return "qb_account_types"
}
