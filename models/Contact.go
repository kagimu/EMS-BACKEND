package models

import (
	"time"
)

type (
	Contact struct {
		ID          uint   `gorm:"primaryKey"`
		Contact     string `json:"contact" form:"contact" binding:"required" gorm:"size:30;unique;not null"`
		ContactType string `json:"contact_type" gorm:"size:10;not null"`
		UserNumber  string `json:"user_number" gorm:"size:100;index;not null"`
		OldNum      string `json:"old_num" gorm:"size:100;"`
		CreatedAt   time.Time
		UpdatedAt   time.Time
		DeletedAt   *time.Time `sql:"index"`
		User        User       `gorm:"foreignKey:UserNumber;references:Number;"`
	}
	ContactList struct {
		ID          uint   `gorm:"primary_key"`
		Contact     string `json:"contact"`
		ContactType string `json:"contact_type"`
		UserNumber  string `json:"user_number"`
	}
)

func (b *Contact) TableName() string {
	return "contacts"
}
