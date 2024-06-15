package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	Number     string         `json:"number" gorm:"size:36;not null;unique"`
	Name       string         `json:"name" gorm:"size:36;not null"`
	Code       int64          `json:"code" gorm:"size:36;not null"`
	UserNumber string         `json:"user_number" gorm:"size:36;not null"`
	User       authModel.User `gorm:"foreignKey:UserNumber;references:Number"`
}
