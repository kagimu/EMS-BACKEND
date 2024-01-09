package models

import (
	"gorm.io/gorm"
)
import authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type MaterialSubType struct {
	gorm.Model
	Number     string         `json:"number" gorm:"size:36;unique;not null;index"`
	Name       string         `json:"name" form:"name" binding:"required" gorm:"size:80;not null"`
	Category   string         `json:"category" form:"category" binding:"required" gorm:"size:80;not null"`
	Code       int64          `json:"code" form:"code" binding:"required" gorm:"size:80;not null"`
	UserNumber string         `json:"user_number" gorm:"size:36;not null"`
	User       authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
