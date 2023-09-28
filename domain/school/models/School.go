package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type (
	School struct {
		ID         uint        `gorm:"primaryKey"`
		Number     string      `json:"number" gorm:"size:100;unique;not null"`
		Name       string      `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		UserNumber string      `json:"user_number" gorm:"size:40;not null"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		DeletedAt  *time.Time  `sql:"index"`
		User       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *School) TableName() string {
	return "schools"
}
