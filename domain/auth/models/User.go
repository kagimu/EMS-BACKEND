package models

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		ID         uint      `gorm:"primaryKey"`
		Number     string    `json:"number" gorm:"size:36;unique;not null"`
		FirstName  string    `json:"first_name" form:"name" binding:"required" gorm:"size:60;not null"`
		SecondName string    `json:"second_name" form:"name" binding:"required" gorm:"size:60;not null"`
		LastName   string    `json:"last_name" form:"name" binding:"required" gorm:"size:60;"`
		Gender     string    `json:"gender" form:"gender" binding:"required" gorm:"size:8;not null"`
		Type       string    `json:"type" gorm:"size:8;not null"`
		Status     string    `json:"status" gorm:"size:100;not null"`
		Contacts   []Contact `json:"contacts" gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *User) TableName() string {
	return "users"
}
