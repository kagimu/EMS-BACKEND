package models

import (
	"gorm.io/gorm"
)

type (
	QBCustomerMapping struct {
		gorm.Model
		UserNumber     string `json:"user_number" gorm:"size:36;uniqueIndex;not null"`
		QbListID       string `json:"qb_list_id" gorm:"size:36;uniqueIndex;not null"`
		QbEditSequence string `json:"qb_edit_sequence" gorm:"size:36"`
		QbFullName     string `json:"qb_full_name" gorm:"size:255"`
	}
)

func (u *QBCustomerMapping) TableName() string {
	return "qb_customer_mapping"
}
