package models

import "time"

type (
	PasswordReset struct {
		ID         uint   `gorm:"primaryKey"`
		Number     string `json:"number" gorm:"size:40;unique;not null"`
		UserNumber string `json:"user_number" gorm:"size:100;index;not null"`
		Token      string `json:"token" gorm:"size:80;not null"`
		Status     string `json:"status" gorm:"size:10;not null"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  *time.Time `sql:"index"`
		User       User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
	PasswordResetInfo struct {
		ID         uint   `gorm:"primary_key"`
		Number     string `json:"number"`
		UserNumber string `json:"user_number"`
		Token      string `json:"token"`
		Status     string `json:"status"`
	}
)

func (b PasswordReset) TableName() string {
	return "password_resets"
}
