package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type Scale struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;not null;unique"`
	Name         string             `json:"name" gorm:"size:1000;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Count        int64              `json:"count"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
