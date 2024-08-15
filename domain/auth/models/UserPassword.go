package models

import "time"

type (
	UserPassword struct {
		ID           uint   `gorm:"primaryKey"`
		UserNumber   string `json:"user_number" gorm:"size:36;unique;index;not null"`
		UserPassword string `json:"user_password" form:"user_password" binding:"required" gorm:"size:70;not null"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    *time.Time `sql:"index"`
		User         User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
	UserPasswordList struct {
		ID           uint   `gorm:"primary_key"`
		UserNumber   string `json:"user_number"`
		UserPassword string `json:"user_password"`
	}
)

func (b UserPassword) TableName() string {
	return "user_passwords"
}
