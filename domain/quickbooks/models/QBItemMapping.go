package models

import (
	"gorm.io/gorm"
)

type (
	QBItemMapping struct {
		gorm.Model
		ItemNumber      string `json:"item_number" gorm:"size:36;not null;uniqueIndex:idx_item_mapping_item_profile"`
		QBProfileNumber string `json:"qb_profile_number" gorm:"size:36;uniqueIndex:idx_item_mapping_item_profile"`
		QbListID        string `json:"qb_list_id" gorm:"size:36;index;not null"`
		QbEditSequence  string `json:"qb_edit_sequence" gorm:"size:36"`
		QbFullName      string `json:"qb_full_name" gorm:"size:255"`
	}
)

func (u *QBItemMapping) TableName() string {
	return "qb_item_mapping"
}
