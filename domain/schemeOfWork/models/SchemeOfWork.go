package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type (
	SchemeOfWork struct {
		gorm.Model
		ID          uint              `gorm:"primaryKey"`
		Number      string            `json:"number" gorm:"size:36;unique;not null"`
		Name        string            `json:"name" form:"name" binding:"required" gorm:"size:200;unique;not null"`
		ClassNumber string            `json:"class_number" gorm:"size:36;not null"`
		UserNumber  string            `json:"user_number" gorm:"size:36;not null"`
		Class       classModels.Class `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User        models.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
